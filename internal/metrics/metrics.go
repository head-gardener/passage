package metrics

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/head-gardener/passage/internal/config"
)

type Metrics struct {
	TunnelStatus *prometheus.GaugeVec
}

func New() (m *Metrics) {
	m = new(Metrics)

	m.TunnelStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "passage",
			Name:      "tunnel_status",
			Help:      "Tunnel status",
		},
		[]string{"remote"},
	)

	return
}

func Serve(log *slog.Logger, conf *config.Config) {
	http.Handle("/metrics", promhttp.Handler())
	log.Debug("serving metrics", "addr", conf.Metrics.Addr)
	http.ListenAndServe(conf.Metrics.Addr, nil)
}
