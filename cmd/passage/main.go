package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/head-gardener/passage/internal/config"
	"github.com/head-gardener/passage/internal/metrics"
	"github.com/head-gardener/passage/internal/net"
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
		log.Error("reading config", "err", err)
		os.Exit(1)
	}

	lvl.Set(conf.Log.Level)

	log.Debug("final config", "val", fmt.Sprintf("%+v", conf))

	var m *metrics.Metrics = nil;

	if conf.Metrics.Enabled {
		m = metrics.New();
		go metrics.Serve(log, conf)
	}

	st, err := net.Init(log, conf, m)
	if err != nil {
		log.Error("initializing", "err", err)
		os.Exit(1)
	}

	go net.Listen(st)
	net.Sender(net.HandleConnection, st)
}
