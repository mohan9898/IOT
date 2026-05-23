package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP 指标
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// 设备指标
	DevicesTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "devices_total",
			Help: "Total number of devices",
		},
	)

	DevicesOnline = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "devices_online",
			Help: "Number of online devices",
		},
	)

	DevicesOffline = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "devices_offline",
			Help: "Number of offline devices",
		},
	)

	// MQTT 指标
	MQTTMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mqtt_messages_total",
			Help: "Total number of MQTT messages",
		},
		[]string{"topic", "direction"},
	)

	MQTTConnectionStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "mqtt_connection_status",
			Help: "MQTT connection status (1=connected, 0=disconnected)",
		},
	)

	// 控制命令指标
	ControlCommandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "control_commands_total",
			Help: "Total number of control commands sent",
		},
		[]string{"device_id", "command"},
	)

	// 认证指标
	LoginAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_attempts_total",
			Help: "Total number of login attempts",
		},
		[]string{"status"},
	)

	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of active users",
		},
	)
)
