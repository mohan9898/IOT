package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusHandler 返回 prometheus HTTP handler
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// MetricsMiddleware 收集 HTTP 指标的中间件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		// 继续处理请求
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		// 记录指标
		HTTPRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		HTTPRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

// RecordMQTTMessage 记录 MQTT 消息
func RecordMQTTMessage(topic string, direction string) {
	MQTTMessagesTotal.WithLabelValues(topic, direction).Inc()
}

// RecordControlCommand 记录控制命令
func RecordControlCommand(deviceID string, command string) {
	ControlCommandsTotal.WithLabelValues(deviceID, command).Inc()
}

// RecordLoginAttempt 记录登录尝试
func RecordLoginAttempt(status string) {
	LoginAttemptsTotal.WithLabelValues(status).Inc()
}

// UpdateDeviceMetrics 更新设备指标
func UpdateDeviceMetrics(total, online, offline int) {
	DevicesTotal.Set(float64(total))
	DevicesOnline.Set(float64(online))
	DevicesOffline.Set(float64(offline))
}

// UpdateActiveUsers 更新活跃用户数
func UpdateActiveUsers(count int) {
	ActiveUsers.Set(float64(count))
}

// SetMQTTConnected 设置 MQTT 连接状态
func SetMQTTConnected(connected bool) {
	if connected {
		MQTTConnectionStatus.Set(1)
	} else {
		MQTTConnectionStatus.Set(0)
	}
}
