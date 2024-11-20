package net

import (
	"github.com/head-gardener/passage/internal/config"
)

type ConnectionHandler func(
	id int,
	done <-chan struct{},
	dev *Device,
	_ *config.Config,
	conn *Connection,
)

func HandleConnection(
	id int,
	done <-chan struct{},
	dev *Device,
	_ *config.Config,
	conn *Connection,
) {
	remote := conn.String()
	dev.Log.Debug("handler started", "remote", remote)
	bufs := make([][]byte, 1024)
	for i := range bufs {
		bufs[i] = make([]byte, 2000)
	}

	for {
		select {
		case <-done:
			dev.Log.Debug("closing connection handler", "remote", remote)
			return
		default:
		}
		n, err := conn.Read(bufs[0])
		if err != nil {
			dev.Log.Error("error reading from peer", "err", err, "remote", remote)
			dev.Close(id)
			continue
		}
		dev.Log.Debug("peer read", "n", n, "remote", remote)

		n, err = dev.Tun.Dev.Write(bufs[0][:n])
		if err != nil {
			dev.Log.Error("error writing to tun", "err", err, "remote", remote)
			continue
		}
		dev.Log.Debug("tun write", "n", n, "remote", remote)
	}
}
