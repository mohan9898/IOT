package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
	"github.com/example/iot-manager/config"
	"github.com/example/iot-manager/internal/metrics"
)

func Connect(cfg *config.MQTTConfig) (mqttlib.Client, error) {
	addr := buildAddr(cfg)
	clientID := fmt.Sprintf("%s%d", cfg.ClientIDPrefix, rand.Int())

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
		log.Println("MQTT connected:", addr)
		metrics.SetMQTTConnected(true)
	}

	opts.OnConnectionLost = func(c mqttlib.Client, err error) {
		log.Println("MQTT connection lost:", err)
		metrics.SetMQTTConnected(false)
	}

	opts.OnReconnecting = func(c mqttlib.Client, opts *mqttlib.ClientOptions) {
		log.Println("MQTT reconnecting...")
	}

	client := mqttlib.NewClient(opts)
	token := client.Connect()

	if token.WaitTimeout(10*time.Second) && token.Error() != nil {
		metrics.SetMQTTConnected(false)
		return nil, token.Error()
	}

	metrics.SetMQTTConnected(true)
	return client, nil
}

func setupTLS(cfg *config.MQTTConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
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
