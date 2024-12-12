package net

import (
	"io"
	"log/slog"
	"net"

	"github.com/head-gardener/passage/internal/config"
	"github.com/net-byte/water"
	"github.com/vishvananda/netlink"
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

	return
}
