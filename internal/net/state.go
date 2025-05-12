package net

import (
	"io"
	"log/slog"
	"net"

	"github.com/net-byte/water"
	"github.com/vishvananda/netlink"

	"github.com/head-gardener/passage/internal/config"
	"github.com/head-gardener/passage/internal/metrics"
)

type State struct {
	conf     *config.Config
	listener net.Listener
	log      *slog.Logger
	metrics  *metrics.Metrics
	netw     *Network
	tun      io.ReadWriteCloser
}

func Init(log *slog.Logger, conf *config.Config, m *metrics.Metrics) (st *State, err error) {
	st = new(State)

	st.listener, err = net.ListenTCP("tcp", &conf.Listener.Addr)
	if err != nil {
		return
	}

	st.netw = New(conf, m)

	devconf := water.Config{DeviceType: water.TUN}
	devconf.Name = conf.Device.Name
	st.tun, err = water.New(devconf)
	if err != nil {
		return
	}

	link, err := netlink.LinkByName(conf.Device.Name)
	if err != nil {
		return
	}
	addr, err := netlink.ParseAddr(conf.Device.Addr)
	if err != nil {
		return
	}
	err = netlink.AddrAdd(link, addr)
	if err != nil {
		return
	}
	err = netlink.LinkSetUp(link)
	if err != nil {
		return
	}

	st.log = log
	st.conf = conf
	st.metrics = m

	return
}
