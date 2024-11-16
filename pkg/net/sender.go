package net

import (
	"github.com/head-gardener/passage/pkg/config"
)

func Send(dev *Device, conf *config.Config) {
	bufs := make([][]byte, 1024)
	for i := range bufs {
		bufs[i] = make([]byte, 100)
	}

	for {
		var err error

		n, err := dev.Tun.Dev.Read(bufs[0])
		if err != nil {
			dev.Log.Error("error reading tun", "err", err)
			continue
		}
		dev.Log.Debug("tun read", "n", n)

		for i := range conf.Peers {
			init, err := dev.EnsureOpen(i, conf)
			if err != nil {
				dev.Log.Error("error connecting to peer", "err", err)
				continue
			}
			if init {
				dev.Log.Info("dialed peer", "peer", conf.Peers[i])
			}

			_, err = dev.Peers[i].conn.Write(bufs[0][:n])
			if err != nil {
				dev.Log.Error("error sending to peer", "err", err)
				dev.Close(i)
				continue
			}
			dev.Log.Debug("peer write", "peer", conf.Peers[i])
		}
	}
}
