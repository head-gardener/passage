package net

import (
	"time"

	"github.com/head-gardener/passage/internal/config"
	"github.com/head-gardener/passage/internal/metrics"
)

type Network struct {
	Peers   []Peer
	metrics *metrics.Metrics
}

type Peer struct {
	conn Connection
	done chan<- struct{}
}

func New(
	conf *config.Config,
	metrcics *metrics.Metrics,
) (netw *Network) {
	netw = new(Network)
	netw.Peers = make([]Peer, len(conf.Peers))
	netw.metrics = metrcics
	return
}

func (netw *Network) EnsureOpen(
	i int,
	handler ConnectionHandler,
	st *State,
) (init bool, err error) {
	init = false
	if netw.IsOpen(i) {
		return
	}

	pass := st.conf.GetSecret()

	err = netw.Peers[i].conn.Dial(&st.conf.Peers[i].Addr, pass)
	if err != nil {
		return
	}
	netw.startConnectionHandler(i, handler, st)

	init = true
	return
}

func (netw *Network) IsOpen(i int) (open bool) {
	return netw.Peers[i].conn.IsOpen()
}

func (netw *Network) sendDone(
	i int,
	st *State,
) {
	select {
	case netw.Peers[i].done <- struct{}{}:
	case <-time.After(3 * time.Second):
		st.log.Error("couldn't close hanging channel, a handler might still be running", "remote", netw.Peers[i].conn.String())
	}
	close(netw.Peers[i].done)
	netw.Peers[i].done = nil
}

func (netw *Network) startConnectionHandler(
	i int,
	handler ConnectionHandler,
	st *State,
) {
	if netw.metrics != nil {
		netw.metrics.TunnelStatus.WithLabelValues(st.conf.Peers[i].Addr.String()).Set(1)
		netw.metrics.ConnectionsEstablished.WithLabelValues(st.conf.Peers[i].Addr.String()).Inc()
	}
	done := make(chan struct{}, 2)
	if netw.Peers[i].done != nil {
		st.log.Warn("hanging channel found, attempting to close...", "remote", netw.Peers[i].conn.String())
		netw.sendDone(i, st)
	}
	netw.Peers[i].done = done
	go handler(i, done, &netw.Peers[i].conn, st)
}

func (netw *Network) Close(i int, st *State) error {
	netw.sendDone(i, st)

	st.log.Info("closing connection", "remote", netw.Peers[i].conn.String())
	closed, err := netw.Peers[i].conn.Close()
	if !closed {
		st.log.Warn("close requested on a nil connection", "remote", netw.Peers[i].conn.String())
	}
	if err != nil && netw.metrics != nil {
		netw.metrics.TunnelStatus.WithLabelValues(st.conf.Peers[i].Addr.String()).Set(0)
	}
	return err
}
