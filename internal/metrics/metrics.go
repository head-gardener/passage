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
	TunnelStatus           *prometheus.GaugeVec
	ConnectionsEstablished *prometheus.CounterVec
	PacketsSent            *prometheus.CounterVec
	PacketsFailed          *prometheus.CounterVec
	PacketsReceived        *prometheus.CounterVec
	PacketsSentOs          prometheus.Counter
	PacketsReceivedOs      prometheus.Counter
	BytesSent              *prometheus.CounterVec
	BytesReceived          *prometheus.CounterVec
	BytesSentOs            prometheus.Counter
	BytesReceivedOs        prometheus.Counter
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

	m.ConnectionsEstablished = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "connections_total",
			Help:      "Total connections established",
		},
		[]string{"remote"},
	)

	m.PacketsSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "packets_sent",
			Help:      "Total packets sent",
		},
		[]string{"remote"},
	)

	m.PacketsFailed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "packets_failed",
			Help:      "Total packets failed to send",
		},
		[]string{"remote"},
	)

	m.PacketsReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "packets_received",
			Help:      "Total packets received",
		},
		[]string{"remote"},
	)

	m.PacketsSentOs = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "packets_sent_os",
			Help:      "Total packets sent to OS",
		},
	)

	m.PacketsReceivedOs = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "packets_received_os",
			Help:      "Total packets received from OS",
		},
	)

	m.BytesSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "bytes_sent",
			Help:      "Total Bytes sent",
		},
		[]string{"remote"},
	)

	m.BytesReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "bytes_received",
			Help:      "Total bytes received",
		},
		[]string{"remote"},
	)

	m.BytesSentOs = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "bytes_sent_os",
			Help:      "Total bytes sent to OS",
		},
	)

	m.BytesReceivedOs = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "passage",
			Name:      "bytes_received_os",
			Help:      "Total bytes received from OS",
		},
	)

	return
}

func Serve(log *slog.Logger, conf *config.Config) {
	http.Handle("/metrics", promhttp.Handler())
	log.Debug("serving metrics", "addr", conf.Metrics.Addr)
	http.ListenAndServe(conf.Metrics.Addr, nil)
}
