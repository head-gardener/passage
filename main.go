package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/head-gardener/passage/config"

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

func listen(dev *Device, conf *config.Config) {
	dev.log.Info("listen init")
	addr, err := net.ResolveUDPAddr("udp", conf.Listener.Addr)
	if err != nil {
		dev.log.Error("error resolving address", "err", err)
		return
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error starting listener", "err", err)
		return
	}
	defer listener.Close()

	for {
		bufs := make([][]byte, 1024)
		for i := range bufs {
			bufs[i] = make([]byte, 100)
		}

		n, addr, err := listener.ReadFromUDP(bufs[0])
		if err != nil {
			dev.log.Error("error reading", "err", err)
			continue
		}
		dev.log.Debug("received", "n", n, "addr", addr, "buf", fmt.Sprintf("% x", bufs[0][:n]))

		n, err = dev.tun.dev.Write(bufs[0][:n])
		if err != nil {
			dev.log.Error("error writing", "err", err)
			continue
		}
		dev.log.Debug("written", "n", n)
	}
}

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))

	conf, err := config.ReadConfig()
	if err != nil {
		log.Error("error reading config", "err", err)
		os.Exit(1)
	}

	lvl.Set(conf.Log.Level)

	log.Debug("final config", "val", fmt.Sprintf("%+v", conf))

	devconf := water.Config{
		DeviceType: water.TUN,
	}
	devconf.Name = conf.Device.Name

	tdev, err := water.New(devconf)
	if err != nil {
		log.Error("error creating tun", "err", err)
		os.Exit(1)
	}
	defer func() {
		tdev.Close()
		log.Debug("tun closed")
	}()

	batchSize := 5
	log.Debug("tun initialised")

	bufs := make([][]byte, batchSize)

	for i := range bufs {
		bufs[i] = make([]byte, 100)
	}

	dev := new(Device)
	dev.log = log
	dev.tun.dev = tdev

	go listen(dev, &conf)

	for {
		var err error

		n, err := tdev.Read(bufs[0])
		if err != nil {
			log.Error("error reading tun device", "err", err)
			continue
		}
		log.Debug("read", "n", n)
		log.Debug("got data", "data", fmt.Sprintf("% x", bufs[0][:n]))

		for i := range conf.Peers {
			conn, err := net.DialUDP("udp", nil, &conf.Peers[i].Addr)
			if err != nil {
				log.Error("error connecting to peer", err, "err")
				continue
			}
			defer conn.Close()

			_, err = conn.Write(bufs[0][:n])
			if err != nil {
				log.Error("error sending to peer", err, "err")
				continue
			}
		}
	}
}
