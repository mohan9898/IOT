package config

import (
	"os"
	"strconv"
	"strings"
)

// Config 是整个 IoT 管理系统的主配置结构体
// 包含了服务器、MQTT、数据库和 JWT 等所有模块的配置
type Config struct {
	Server   ServerConfig   // HTTP 服务器配置
	MQTT     MQTTConfig     // MQTT 消息队列配置
	Database DatabaseConfig // 数据库配置
	JWT      JWTConfig      // JWT 认证配置
}

// ServerConfig 定义了 HTTP 服务器的配置参数
type ServerConfig struct {
	Host               string   // 服务器监听的地址，0.0.0.0 表示监听所有网络接口
	Port               int      // 服务器监听的端口号
	HTTPS              bool     // 是否启用 HTTPS
	CertPath           string   // HTTPS 证书文件路径
	KeyPath            string   // HTTPS 私钥文件路径
	CORSAllowedOrigins []string // CORS 允许的源列表
}

// MQTTConfig 定义了 MQTT 消息代理的连接配置
// 用于与 IoT 设备进行通信
type MQTTConfig struct {
	Protocol       string // 连接协议，支持 tcp 或 ssl/tls
	Broker         string // MQTT 代理服务器的地址
	Port           int    // MQTT 代理服务器的端口
	Path           string // MQTT 连接路径（主要用于 WebSocket 连接）
	ClientIDPrefix string // MQTT 客户端 ID 的前缀
	Username       string // MQTT 连接用户名（如需要认证）
	Password       string // MQTT 连接密码（如需要认证）
	TLSEnabled     bool   // 是否启用 TLS 加密连接
	CACertPath     string // CA 证书文件路径（启用 TLS 时使用）
}

// DatabaseConfig 定义了数据库的配置
// 当前系统使用 SQLite 嵌入式数据库
type DatabaseConfig struct {
	Path         string // SQLite 数据库文件的存储路径
	BackupPath   string // 数据库备份路径
	BackupEnable bool   // 是否启用自动备份
	BackupHours  int    // 备份间隔（小时）
}

// JWTConfig 定义了 JWT（JSON Web Token）认证的配置
// 用于用户身份验证和会话管理
type JWTConfig struct {
	Secret       string // JWT 签名密钥，生产环境必须修改为强随机字符串
	ExpiresHours int    // JWT Token 的有效期，单位为小时
}

// Load 函数加载并返回系统配置
// 首先使用默认值初始化配置，然后可以通过环境变量覆盖部分配置
// 返回值：指向完整配置对象的指针
func Load() *Config {
	// 初始化默认配置
	cfg := &Config{
		Server: ServerConfig{
			Host:               "0.0.0.0",            // 监听所有网络接口
			Port:               6116,                 // 默认端口号 6116
			HTTPS:              false,                // 默认不启用 HTTPS
			CertPath:           "",                   // 证书路径
			KeyPath:            "",                   // 私钥路径
			CORSAllowedOrigins: []string{},           // 默认允许的源（空表示使用默认策略）
		},
		MQTT: MQTTConfig{
			Protocol: "ssl",                                     // 使用 TLS/SSL 协议
			Broker: "d11aab19.ala.cn-hangzhou.emqxsl.cn",                     // 新的 EMQX 服务器地址
			Port: 8883,                                      // TLS/SSL 端口
			Path: "/mqtt",                                     // WebSocket 路径
			ClientIDPrefix: "iot-manager-",                            // 客户端 ID 前缀
			Username: "taiyi",                        // MQTT 用户名（通过环境变量配置）
Password: "p9UPpz2i.H48stE",             // MQTT 密码（通过环境变量配置）
			TLSEnabled: true,                                     // 启用 TLS 加密连接
			CACertPath: "",                                        // CA 证书路径（使用系统根证书）
		},
		Database: DatabaseConfig{
			Path: "./data/iot.db", // 数据库文件存储在 data 目录下
			BackupPath: "./data/backup", // 备份路径
			BackupEnable: false,         // 默认不启用自动备份
			BackupHours: 24,             // 每 24 小时备份一次
		},
		JWT: JWTConfig{
			Secret: "", // JWT 密钥（通过环境变量配置，必须设置）
			ExpiresHours: 24, // Token 有效期为 24 小时
		},
	}

	// 通过环境变量覆盖配置
	if os.Getenv("SERVER_HOST") != "" {
		cfg.Server.Host = os.Getenv("SERVER_HOST")
	}
	if os.Getenv("SERVER_PORT") != "" {
		if port, err := strconv.Atoi(os.Getenv("SERVER_PORT")); err == nil {
			cfg.Server.Port = port
		}
	}
	if os.Getenv("HTTPS_ENABLE") != "" {
		if https, err := strconv.ParseBool(os.Getenv("HTTPS_ENABLE")); err == nil {
			cfg.Server.HTTPS = https
		}
	}
	if os.Getenv("HTTPS_CERT") != "" {
		cfg.Server.CertPath = os.Getenv("HTTPS_CERT")
	}
	if os.Getenv("HTTPS_KEY") != "" {
		cfg.Server.KeyPath = os.Getenv("HTTPS_KEY")
	}

	// 通过环境变量覆盖 MQTT 配置
	if os.Getenv("MQTT_BROKER") != "" {
		cfg.MQTT.Broker = os.Getenv("MQTT_BROKER")
	}
	if os.Getenv("MQTT_PORT") != "" {
		if port, err := strconv.Atoi(os.Getenv("MQTT_PORT")); err == nil {
			cfg.MQTT.Port = port
		}
	}
	if os.Getenv("MQTT_USERNAME") != "" {
		cfg.MQTT.Username = os.Getenv("MQTT_USERNAME")
	}
	if os.Getenv("MQTT_PASSWORD") != "" {
		cfg.MQTT.Password = os.Getenv("MQTT_PASSWORD")
	}
	if os.Getenv("MQTT_PROTOCOL") != "" {
		cfg.MQTT.Protocol = os.Getenv("MQTT_PROTOCOL")
	}
	if os.Getenv("MQTT_TLS_ENABLED") != "" {
		if enabled, err := strconv.ParseBool(os.Getenv("MQTT_TLS_ENABLED")); err == nil {
			cfg.MQTT.TLSEnabled = enabled
		}
	}

	// 数据库配置
	if os.Getenv("DB_PATH") != "" {
		cfg.Database.Path = os.Getenv("DB_PATH")
	}
	if os.Getenv("DB_BACKUP_ENABLE") != "" {
		if enable, err := strconv.ParseBool(os.Getenv("DB_BACKUP_ENABLE")); err == nil {
			cfg.Database.BackupEnable = enable
		}
	}
	if os.Getenv("DB_BACKUP_PATH") != "" {
		cfg.Database.BackupPath = os.Getenv("DB_BACKUP_PATH")
	}
	if os.Getenv("DB_BACKUP_HOURS") != "" {
		if hours, err := strconv.Atoi(os.Getenv("DB_BACKUP_HOURS")); err == nil {
			cfg.Database.BackupHours = hours
		}
	}

	// JWT 配置
	if os.Getenv("JWT_SECRET") != "" {
		cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	}
	if os.Getenv("JWT_EXPIRES_HOURS") != "" {
		if hours, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_HOURS")); err == nil {
			cfg.JWT.ExpiresHours = hours
		}
	}

	// CORS 配置
	if os.Getenv("CORS_ALLOWED_ORIGINS") != "" {
		origins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				cfg.Server.CORSAllowedOrigins = append(cfg.Server.CORSAllowedOrigins, origin)
			}
		}
	}

	return cfg
}
