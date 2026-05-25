package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	mqttlib "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/example/iot-manager/config"
	"github.com/example/iot-manager/internal/api"
	"github.com/example/iot-manager/internal/db"
	"github.com/example/iot-manager/internal/mqtt"
)

func main() {
	// 设置日志
	logger, err := setupLogger()
	if err != nil {
		fmt.Printf("Failed to setup logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 加载配置
	cfg := config.Load()
	logger.Info("Configuration loaded",
		zap.String("server_host", cfg.Server.Host),
		zap.Int("server_port", cfg.Server.Port),
		zap.Bool("https_enabled", cfg.Server.HTTPS),
		zap.String("mqtt_broker", cfg.MQTT.Broker),
		zap.String("db_path", cfg.Database.Path),
		zap.Bool("backup_enabled", cfg.Database.BackupEnable),
	)

	// 初始化数据库
	database, err := db.NewSQLite(cfg.Database.Path)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.Close()
	logger.Info("Database initialized successfully")

	// 启动备份管理器
	var backupManager *db.BackupManager
	if cfg.Database.BackupEnable {
		backupManager = db.NewBackupManager(
			database,
			cfg.Database.Path,
			cfg.Database.BackupPath,
			cfg.Database.BackupHours,
			logger,
		)
		backupManager.Start()
		defer backupManager.Stop()
		logger.Info("Backup manager started")
	}

	// 连接 MQTT
	mqttClient, err := mqtt.Connect(&cfg.MQTT)
	if err != nil {
		logger.Warn("Failed to connect to MQTT broker, continuing without MQTT", zap.Error(err))
		// 创建一个 dummy 客户端，这样应用程序不会崩溃
		mqttClient = nil
	} else {
		logger.Info("MQTT connected successfully")
	}

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(ginLogger(logger))

	// 创建 API 处理器
	handler := api.NewHandler(database, mqttClient, cfg, logger)
	handler.SetupRoutes(r)

	// 订阅 MQTT 主题
	if mqttClient != nil {
		topics := []string{
			"smart_light/#",
			"+/register",
			"+/status",
			"+/control",
			"+/metric",
			"$SYS/brokers/+/clients/+/connected",
			"$SYS/brokers/+/clients/+/disconnected",
		}

		for _, topic := range topics {
			token := mqttClient.Subscribe(topic, 0, func(client mqttlib.Client, msg mqttlib.Message) {
				handler.HandleMQTTMessage(msg)
			})
			token.Wait()
			logger.Info("Subscribed to topic", zap.String("topic", topic))
		}
	}

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Starting server", zap.String("address", serverAddr))

	go func() {
		var err error
		if cfg.Server.HTTPS {
			logger.Info("Starting HTTPS server")
			err = r.RunTLS(serverAddr, cfg.Server.CertPath, cfg.Server.KeyPath)
		} else {
			logger.Info("Starting HTTP server")
			err = r.Run(serverAddr)
		}
		if err != nil {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			count, err := database.MarkOfflineDevices(5)
			if err != nil {
				logger.Warn("Failed to mark offline devices", zap.Error(err))
			} else if count > 0 {
				logger.Info("Devices marked offline", zap.Int("count", count))
			}
		}
	}()

	// 等待信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gracefully...")
	mqttClient.Disconnect(250)
    logger.Info("Server shutdown complete")
}

func setupLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.CallerKey = "caller"
	
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	
	return logger, nil
}

func ginLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		logger.Info("HTTP request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		)
	}
}