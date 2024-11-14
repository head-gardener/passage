package net

import (
	"github.com/head-gardener/passage/pkg/config"
)

func HandleConnection(dev *Device, _ *config.Config, conn *Connection) {
	dev.Log.Debug("handler started", "remote", conn.tcp.RemoteAddr().String())
	for {
		bufs := make([][]byte, 1024)
		for i := range bufs {
			bufs[i] = make([]byte, 100)
		}

		n, err := conn.Read(bufs[0])
		if err != nil {
			dev.Log.Error("error reading from peer", "err", err)
			continue
		}
		dev.Log.Debug("peer read", "n", n)

		n, err = dev.Tun.Dev.Write(bufs[0][:n])
		if err != nil {
			dev.Log.Error("error writing to tun", "err", err)
			continue
		}
		dev.Log.Debug("tun write", "n", n)
	}
}
