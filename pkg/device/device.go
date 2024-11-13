package device

import (
	"log/slog"

	"github.com/head-gardener/passage/pkg/config"
	"github.com/net-byte/water"
)

type Device struct {
	State struct {
	}

	Rate struct {
	}

	Tun struct {
		Dev *water.Interface
	}

	Log *slog.Logger
}

func New(conf *config.Config, log *slog.Logger) (dev *Device, err error) {
	devconf := water.Config{
		DeviceType: water.TUN,
	}
	devconf.Name = conf.Device.Name

	tdev, err := water.New(devconf)

	dev = new(Device)
	dev.Log = log
	dev.Tun.Dev = tdev

	return
}
