package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
	"github.com/example/iot-manager/config"
	"github.com/example/iot-manager/internal/metrics"
)

type ConnectionStatus struct {
	Connected   bool      `json:"connected"`
	Broker      string    `json:"broker"`
	Port        int       `json:"port"`
	Protocol    string    `json:"protocol"`
	TLSEnabled  bool      `json:"tls_enabled"`
	ConnectedAt time.Time `json:"connected_at"`
}

var (
	currentStatus = ConnectionStatus{}
	statusMu      sync.RWMutex
)

func GetStatus() ConnectionStatus {
	statusMu.RLock()
	defer statusMu.RUnlock()
	return currentStatus
}

func Connect(cfg *config.MQTTConfig) (mqttlib.Client, error) {
	addr := buildAddr(cfg)
	clientID := fmt.Sprintf("%s%d", cfg.ClientIDPrefix, rand.Int())

	statusMu.Lock()
	currentStatus = ConnectionStatus{
		Connected:  false,
		Broker:     cfg.Broker,
		Port:       cfg.Port,
		Protocol:   cfg.Protocol,
		TLSEnabled: cfg.TLSEnabled,
	}
	statusMu.Unlock()

	opts := mqttlib.NewClientOptions()
	opts.AddBroker(addr)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	if cfg.TLSEnabled {
		tlsConfig, err := setupTLS(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to setup TLS: %w", err)
		}
		opts.SetTLSConfig(tlsConfig)
	}

	opts.OnConnect = func(c mqttlib.Client) {
		fmt.Printf("[MQTT] 已连接: %s (KeepAlive=%ds)\n", addr, 60)
		metrics.SetMQTTConnected(true)
		statusMu.Lock()
		currentStatus.Connected = true
		currentStatus.ConnectedAt = time.Now()
		statusMu.Unlock()
	}

	opts.OnConnectionLost = func(c mqttlib.Client, err error) {
		fmt.Printf("[MQTT] 连接断开: %v (客户端状态=%v)\n", err, c.IsConnected())
		metrics.SetMQTTConnected(false)
		statusMu.Lock()
		currentStatus.Connected = false
		statusMu.Unlock()
	}

	opts.OnReconnecting = func(c mqttlib.Client, opts *mqttlib.ClientOptions) {
		fmt.Println("[MQTT] 正在重连...")
	}

	client := mqttlib.NewClient(opts)
	token := client.Connect()

	if token.WaitTimeout(10*time.Second) && token.Error() != nil {
		metrics.SetMQTTConnected(false)
		return nil, token.Error()
	}

	metrics.SetMQTTConnected(true)
	statusMu.Lock()
	currentStatus.Connected = true
	currentStatus.ConnectedAt = time.Now()
	statusMu.Unlock()

	return client, nil
}

func setupTLS(cfg *config.MQTTConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.CACertPath == "",
	}

	if cfg.CACertPath != "" {
		caCert, err := os.ReadFile(cfg.CACertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA cert: %w", err)
		}
		
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}

func buildAddr(cfg *config.MQTTConfig) string {
	if cfg.Protocol == "ws" || cfg.Protocol == "wss" {
		return fmt.Sprintf("%s://%s:%d%s", cfg.Protocol, cfg.Broker, cfg.Port, cfg.Path)
	}
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Broker, cfg.Port)
}
