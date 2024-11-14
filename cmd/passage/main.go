package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/head-gardener/passage/pkg/config"
	"github.com/head-gardener/passage/pkg/net"
)

func initLog() (log *slog.Logger, lvl *slog.LevelVar) {
	lvl = new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))

	return
}

func main() {
	log, lvl := initLog()

	conf, err := config.ReadConfig()
	if err != nil {
		log.Error("error reading config", "err", err)
		os.Exit(1)
	}

	lvl.Set(conf.Log.Level)

	log.Debug("final config", "val", fmt.Sprintf("%+v", conf))

	dev, err := net.New(&conf, log, net.HandleConnection)
	if err != nil {
		log.Error("error initializing device", "err", err)
		os.Exit(1)
	}

	go net.Listen(dev, &conf)
	net.Send(dev, &conf)
}
