package net

import (
	"log/slog"

	"github.com/head-gardener/passage/internal/config"
)

type Network struct {
	Peers []Peer
}

type Peer struct {
	conn Connection
	done chan<- struct{}
}

func New(
	conf *config.Config,
) (netw *Network) {
	netw = new(Network)
	netw.Peers = make([]Peer, len(conf.Peers))
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

	pass, err := st.conf.GetSecret()
	if err != nil {
		return
	}

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

func (netw *Network) startConnectionHandler(
	i int,
	handler ConnectionHandler,
	st *State,
) {
	done := make(chan struct{}, 2)
	go handler(i, done, &netw.Peers[i].conn, st)
}

func (netw *Network) Close(i int, log *slog.Logger) error {
	netw.Peers[i].done <- struct{}{}
	log.Info("closing connection", "remote", netw.Peers[i].conn.String())
	closed, err := netw.Peers[i].conn.Close()
	if !closed {
		log.Warn("close requested on a nil connection")
	}
	return err
}
