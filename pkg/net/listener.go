package net

import (
	"fmt"
	"net"

	"github.com/head-gardener/passage/pkg/config"
	"github.com/head-gardener/passage/pkg/device"
)

func Listen(dev *device.Device, conf *config.Config) {
	listener, err := net.ListenTCP("tcp", &conf.Listener.Addr)
	if err != nil {
		fmt.Println("error starting listener", "err", err)
		return
	}
	defer listener.Close()

	dev.Log.Info("listener initialized", "addr", conf.Listener.Addr)

	for {

		conn, err := listener.Accept()
		// TODO: multithread
		if err != nil {
			fmt.Println("error accepting connection", "err", err)
			return
		}
		handleConnection(dev, conf, conn)
	}
}

func handleConnection(dev *device.Device, _ *config.Config, conn net.Conn) {
	dev.Log.Debug("accepted connection")
	for {
		bufs := make([][]byte, 1024)
		for i := range bufs {
			bufs[i] = make([]byte, 100)
		}

		n, err := conn.Read(bufs[0])
		if err != nil {
			dev.Log.Error("error reading", "err", err)
			continue
		}
		dev.Log.Debug("received", "n", n, "addr", conn.RemoteAddr().String(), "buf", fmt.Sprintf("% x", bufs[0][:n]))

		n, err = dev.Tun.Dev.Write(bufs[0][:n])
		if err != nil {
			dev.Log.Error("error writing", "err", err)
			continue
		}
		dev.Log.Debug("written", "n", n)
	}
}
