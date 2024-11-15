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

	handler func(*Device, *config.Config, *Connection)
}

type Peer struct {
	conn Connection
}

func New(
	conf *config.Config,
	log *slog.Logger,
	handler func(*Device, *config.Config, *Connection),
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

func (dev *Device) EnsureOpen(i int, conf *config.Config) (bool, error) {
	if dev.Peers[i].conn.tcp != nil {
		return false, nil
	}

	err := dev.Peers[i].conn.Dial(&conf.Peers[i].Addr)
	if err != nil {
		return false, err
	}
	go dev.handler(dev, conf, &dev.Peers[i].conn)
	return true, nil
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
	for i := range conf.Peers {
		if !conf.Peers[i].Addr.IP.Equal(remoteAddr) {
			dev.Log.Debug("no match", "laddr", conf.Peers[i].Addr.IP.String(), "raddr", remoteAddr.String())
			continue
		}

		found = true

		dev.Log.Info("initialized connection", "remote", remoteAddr.String())
		if dev.Peers[i].conn.tcp != nil {
			dev.Peers[i].conn.tcp.Close()
		}
		dev.Peers[i].conn.tcp = tcp

		go dev.handler(dev, conf, &dev.Peers[i].conn)
	}
	if !found {
		dev.Log.Warn("unexpected connection", "addr", remoteAddr.String())
		tcp.Close()
	}

	return
}
