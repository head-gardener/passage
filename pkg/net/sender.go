package net

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/head-gardener/passage/pkg/config"
	"github.com/head-gardener/passage/pkg/device"

	"github.com/net-byte/water"
)

type Device struct {
	state struct {
	}

	rate struct {
	}

	tun struct {
		dev *water.Interface
	}

	log *slog.Logger
}

func Send(dev *device.Device, conf *config.Config) {
	bufs := make([][]byte, 1024)
	for i := range bufs {
		bufs[i] = make([]byte, 100)
	}

	for {
		var err error

		n, err := dev.Tun.Dev.Read(bufs[0])
		if err != nil {
			dev.Log.Error("error reading tun device", "err", err)
			continue
		}
		dev.Log.Debug("read", "n", n)
		dev.Log.Debug("got data", "data", fmt.Sprintf("% x", bufs[0][:n]))

		for i := range conf.Peers {
			conn, err := net.DialTCP("tcp", nil, &conf.Peers[i].Addr)
			if err != nil {
				dev.Log.Error("error connecting to peer", "err", err)
				continue
			}
			defer conn.Close()

			_, err = conn.Write(bufs[0][:n])
			if err != nil {
				dev.Log.Error("error sending to peer", "err", err)
				continue
			}
			dev.Log.Debug("sent data to peer", "peer", conf.Peers[i])
		}
	}
}
