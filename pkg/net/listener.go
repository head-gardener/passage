package net

import (
	"github.com/head-gardener/passage/pkg/config"
)

func Listen(dev *Device, conf *config.Config) {
	err := dev.InitListener(conf)
	if err != nil {
		dev.Log.Error("error initializing listener", "err", err)
	}

	for {
		dev.Accept(conf)
		if err != nil {
			dev.Log.Error("error initializing listener", "err", err)
		}
	}
}
