package net

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/head-gardener/passage/pkg/config"

	"github.com/net-byte/water"
)

type Device struct {
	Peers []Peer

	Rate struct {
	}

	listener struct {
		tcp *net.TCPListener
	}

	Tun struct {
		Dev *water.Interface
	}

	Log *slog.Logger

	handler ConnectionHandler
}

type Peer struct {
	conn Connection
	done chan<- struct{}
}

func New(
	conf *config.Config,
	log *slog.Logger,
	handler ConnectionHandler,
) (dev *Device, err error) {
	devconf := water.Config{
		DeviceType: water.TUN,
	}
	devconf.Name = conf.Device.Name

	tdev, err := water.New(devconf)

	dev = new(Device)
	dev.Log = log
	dev.Tun.Dev = tdev
	dev.Peers = make([]Peer, len(conf.Peers))
	dev.handler = handler

	return
}

func (dev *Device) EnsureOpen(i int, conf *config.Config) (init bool, err error) {
	init = false
	if dev.Peers[i].conn.tcp != nil {
		return
	}

	pass, err := conf.GetSecret()
	if err != nil {
		return
	}

	err = dev.Peers[i].conn.Dial(&conf.Peers[i].Addr, pass)
	if err != nil {
		return
	}
	dev.startConnectionHandler(i, conf)

	init = true
	return
}

func (dev *Device) startConnectionHandler(i int, conf *config.Config) {
	done := make(chan struct{}, 2)
	dev.Peers[i].done = done
	go dev.handler(i, done, dev, conf, &dev.Peers[i].conn)
}

func (dev *Device) Close(i int) error {
	dev.Peers[i].done <- struct{}{}
	dev.Log.Info("closing connection", "remote", dev.Peers[i].conn.String())
	closed, err := dev.Peers[i].conn.Close()
	if !closed {
		dev.Log.Warn("close requested on a nil connection")
	}
	return err
}

func (dev *Device) InitListener(conf *config.Config) (err error) {
	listener, err := net.ListenTCP("tcp", &conf.Listener.Addr)
	if err != nil {
		return
	}

	dev.listener.tcp = listener
	return
}

func (dev *Device) Accept(conf *config.Config) (err error) {
	tcp, err := dev.listener.tcp.AcceptTCP()
	dev.Log.Debug("received connection")
	if err != nil {
		return fmt.Errorf("couldn't initialize tcp listener: %w", err)
	}

	remote, ok := tcp.RemoteAddr().(*net.TCPAddr)
	if !ok {
		tcp.Close()
		return fmt.Errorf("couldn't coerce remote addr to TCPAddr: %v", tcp.RemoteAddr())
	}
	remoteAddr := remote.IP

	found := false

	i := 0
	for ; i < len(conf.Peers); i++ {
		if conf.Peers[i].Addr.IP.Equal(remoteAddr) {
			found = true
			break
		}
	}

	if !found {
		tcp.Close()
		return fmt.Errorf("unexpected connection from %s", remoteAddr.String())
	}

	dev.Log.Info("initialized connection", "remote", remoteAddr.String())
	if dev.Peers[i].conn.tcp != nil {
		tcp.Close()
		return fmt.Errorf("peer already established connection")
	}

	pass, err := conf.GetSecret()
	if err != nil {
		tcp.Close()
		return err
	}

	err = dev.Peers[i].conn.Accept(tcp, pass)
	if err != nil {
		tcp.Close()
		return err
	}

	dev.startConnectionHandler(i, conf)

	return
}
