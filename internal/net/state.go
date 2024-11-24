package net

import (
	"io"
	"log/slog"
	"net"

	"github.com/head-gardener/passage/internal/config"
	"github.com/net-byte/water"
)

type State struct {
	listener net.Listener
	netw     *Network
	tun      io.ReadWriteCloser
	log      *slog.Logger
	conf     *config.Config
}

func Init(log *slog.Logger, conf *config.Config) (st *State, err error) {
	st = new(State)

	st.listener, err = net.ListenTCP("tcp", &conf.Listener.Addr)
	if err != nil {
		return
	}

	st.netw = New(conf)

	devconf := water.Config{DeviceType: water.TUN}
	devconf.Name = conf.Device.Name
	st.tun, err = water.New(devconf)
	if err != nil {
		return
	}

	st.log = log
	st.conf = conf

	return
}
