package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/example/iot-manager/config"
	"github.com/example/iot-manager/internal/auth"
	"github.com/example/iot-manager/internal/db"
	"github.com/example/iot-manager/internal/metrics"
	mqttstatus "github.com/example/iot-manager/internal/mqtt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db                 *db.SQLite
	mqtt               mqtt.Client
	wsHub              *WebSocketHub
	jwtManager         *auth.JWTManager
	logger             *zap.Logger
	rateLimit          map[string]int
	rateMu             sync.Mutex
	allowedOrigins     []string
}

func NewHandler(db *db.SQLite, mqtt mqtt.Client, cfg *config.Config, logger *zap.Logger) *Handler {
	return &Handler{
		db:             db,
		mqtt:           mqtt,
		wsHub:          NewWebSocketHub(),
		jwtManager:     auth.NewJWTManager(&cfg.JWT),
		logger:         logger,
		rateLimit:      make(map[string]int),
		allowedOrigins: cfg.Server.CORSAllowedOrigins,
	}
}

func (h *Handler) SetupRoutes(r *gin.Engine) {
	go h.wsHub.Run()
	go h.cleanupRateLimit()

	r.Use(h.SecurityHeadersMiddleware())
	r.Use(h.CORSMiddleware())
	r.Use(metrics.MetricsMiddleware())

	// Prometheus 指标
	r.GET("/metrics", metrics.PrometheusHandler())

	// 健康检查
	r.GET("/health", h.HealthCheck)

	// 认证路由（无需认证）
	r.POST("/api/auth/login", h.Login)
	r.POST("/api/auth/register", h.Register)

	// 需要认证的路由
	api := r.Group("/api")
	api.Use(h.RateLimitMiddleware())
	api.Use(h.AuthMiddleware())
	{
		api.POST("/auth/update-account", h.UpdateAccount)
		api.POST("/auth/update-password", h.UpdatePassword)
		api.GET("/auth/user-info", h.GetUserInfo)

		api.GET("/devices", h.GetDevices)
		api.GET("/devices/:id", h.GetDevice)
		api.POST("/devices", h.CreateDevice)
		api.PUT("/devices/:id", h.UpdateDevice)
		api.DELETE("/devices/:id", h.DeleteDevice)
		api.GET("/devices/stats", h.GetDeviceStats)

		api.GET("/device-types", h.GetDeviceTypes)
		api.POST("/device-types", h.AddDeviceType)
		api.GET("/device-types/:type_id", h.GetDeviceType)

		api.POST("/control/send", h.SendCommand)
		api.POST("/control/threshold", h.SetThreshold)
		api.GET("/control/history/:device_id", h.GetCommandHistory)

		api.GET("/control/records", h.GetControlRecords)
		api.GET("/control/stats", h.GetControlStats)

		api.GET("/ws", h.WebSocketHandler)

		api.GET("/dashboard", h.GetDashboard)

		api.GET("/mqtt/status", h.GetMQTTStatus)
	}

	// 静态文件服务 — 路径可通过 STATIC_DIR 环境变量自定义
	// Docker 中默认为 ./dist，本地开发可设置 STATIC_DIR=../dist
	staticDir := "./dist"
	if d := os.Getenv("STATIC_DIR"); d != "" {
		staticDir = d
	}
	r.Static("/web", staticDir)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/")
	})
}

// RateLimitMiddleware 速率限制中间件
func (h *Handler) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		h.rateMu.Lock()
		count := h.rateLimit[clientIP]
		h.rateLimit[clientIP] = count + 1
		h.rateMu.Unlock()

		if count > 100 { // 每分钟100次请求限制
			h.logger.Warn("Rate limit exceeded", zap.String("ip", clientIP))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "请求过于频繁，请稍后再试"})
			return
		}
		c.Next()
	}
}

// cleanupRateLimit 定期清理速率限制记录
func (h *Handler) cleanupRateLimit() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		h.rateMu.Lock()
		h.rateLimit = make(map[string]int)
		h.rateMu.Unlock()
	}
}

// AuthMiddleware JWT 认证中间件
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}
		
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := h.jwtManager.ValidateToken(token)
		if err != nil {
			h.logger.Warn("Invalid token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		
		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// HealthCheck 健康检查端点
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		metrics.RecordLoginAttempt("invalid_request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.db.GetUserByUsername(req.Username)
	if err != nil {
		h.logger.Warn("Login failed: user not found", zap.String("username", req.Username))
		metrics.RecordLoginAttempt("failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.logger.Warn("Login failed: invalid password", zap.String("username", req.Username))
		metrics.RecordLoginAttempt("failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		metrics.RecordLoginAttempt("failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	h.logger.Info("User logged in", zap.String("username", req.Username))
	metrics.RecordLoginAttempt("success")
	metrics.UpdateActiveUsers(1) // 简化版本
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 3 characters"})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	user, err := h.db.CreateUser(req.Username, string(hash))
	if err != nil {
		h.logger.Warn("Registration failed: username exists", zap.String("username", req.Username))
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	token, err := h.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	h.logger.Info("User registered", zap.String("username", req.Username))
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.Username != "" {
		user, err := h.db.GetUserByUsername(req.Username)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"exists": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"exists":   true,
			"username": user.Username,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Username required"})
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	var req struct {
		OldUsername string `json:"old_username"`
		NewUsername string `json:"new_username"`
		Password    string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.NewUsername) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 3 characters"})
		return
	}

	user, err := h.db.GetUserByUsername(req.OldUsername)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	checkUser, err := h.db.GetUserByUsername(req.NewUsername)
	if err == nil && checkUser != nil && checkUser.ID != user.ID {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	if err := h.db.UpdateUser(user.ID, req.NewUsername, ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	h.logger.Info("Account updated", zap.String("old_username", req.OldUsername), zap.String("new_username", req.NewUsername))
	c.JSON(http.StatusOK, gin.H{"message": "Username updated", "username": req.NewUsername})
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	var req struct {
		Username    string `json:"username"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	user, err := h.db.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	if err := h.db.UpdatePassword(user.ID, string(newHash)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	h.logger.Info("Password updated", zap.String("username", req.Username))
	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}

func (h *Handler) GetDevices(c *gin.Context) {
	devices, err := h.db.GetDevices()
	if err != nil {
		h.logger.Error("Failed to get devices", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *Handler) GetDevice(c *gin.Context) {
	id := c.Param("id")
	device, err := h.db.GetDevice(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (h *Handler) CreateDevice(c *gin.Context) {
	var req struct {
		ID       string                 `json:"id"`
		Name     string                 `json:"name"`
		Type     string                 `json:"type"`
		Metadata map[string]interface{} `json:"metadata,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}

	device := &db.Device{
		ID:        req.ID,
		Name:      req.Name,
		Type:      req.Type,
		Status:    "offline",
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
	}

	if err := h.db.CreateDevice(device); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Device already exists"})
		return
	}

	h.logger.Info("Device created", zap.String("device_id", req.ID), zap.String("device_name", req.Name))
	c.JSON(http.StatusCreated, device)
}

func (h *Handler) UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name     string                 `json:"name"`
		Type     string                 `json:"type,omitempty"`
		Metadata map[string]interface{} `json:"metadata,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	device, err := h.db.GetDevice(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Type != "" {
		device.Type = req.Type
	}
	if req.Metadata != nil {
		device.Metadata = req.Metadata
	}

	if err := h.db.UpdateDevice(device); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Device updated", zap.String("device_id", id))
	c.JSON(http.StatusOK, device)
}

func (h *Handler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.DeleteDevice(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Device deleted", zap.String("device_id", id))
	c.JSON(http.StatusOK, gin.H{"message": "Device deleted"})
}

func (h *Handler) GetDeviceStats(c *gin.Context) {
	total, online, offline, err := h.db.GetDeviceStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 更新 Prometheus 指标
	metrics.UpdateDeviceMetrics(total, online, offline)
	c.JSON(http.StatusOK, gin.H{
		"total":   total,
		"online":  online,
		"offline": offline,
	})
}

func (h *Handler) GetDashboard(c *gin.Context) {
	total, online, offline, _ := h.db.GetDeviceStats()

	devices, _ := h.db.GetDevices()
	typeCount := make(map[string]int)
	var recentDevices []gin.H
	for i, d := range devices {
		typeCount[d.Type]++
		if i >= len(devices)-5 {
			recentDevices = append(recentDevices, gin.H{
				"id":   d.ID,
				"name": d.Name,
				"type": d.Type,
				"status": d.Status,
			})
		}
	}

	var typeDistribution []gin.H
	for typeID, count := range typeCount {
		dt, _ := h.db.GetDeviceType(typeID)
		name := typeID
		icon := "📦"
		if dt != nil {
			if n, ok := dt["name"].(string); ok {
				name = n
			}
			if ic, ok := dt["icon"].(string); ok {
				icon = ic
			}
		}
		typeDistribution = append(typeDistribution, gin.H{
			"type_id": typeID,
			"name":    name,
			"icon":    icon,
			"count":   count,
		})
	}
	if typeDistribution == nil {
		typeDistribution = []gin.H{}
	}
	if recentDevices == nil {
		recentDevices = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"total":   total,
			"online":  online,
			"offline": offline,
		},
		"type_distribution": typeDistribution,
		"recent_devices":    recentDevices,
	})
}

func (h *Handler) GetDeviceTypes(c *gin.Context) {
	types, err := h.db.GetDeviceTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types)
}

func (h *Handler) AddDeviceType(c *gin.Context) {
	var req struct {
		TypeID        string   `json:"type_id"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Icon          string   `json:"icon,omitempty"`
		Commands      []string `json:"commands"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Icon == "" {
		req.Icon = "📦"
	}

	if err := h.db.AddDeviceType(req.TypeID, req.Name, req.Description, req.Icon, req.Commands); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Device type added", zap.String("type_id", req.TypeID))
	c.JSON(http.StatusCreated, gin.H{"message": "Device type added", "type_id": req.TypeID})
}

func (h *Handler) GetDeviceType(c *gin.Context) {
	typeID := c.Param("type_id")
	deviceType, err := h.db.GetDeviceType(typeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device type not found"})
		return
	}
	c.JSON(http.StatusOK, deviceType)
}

func (h *Handler) SendCommand(c *gin.Context) {
	var req struct {
		DeviceID   string                 `json:"device_id"`
		Command    string                 `json:"command"`
		Parameters map[string]interface{} `json:"parameters,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Parameters == nil {
		req.Parameters = make(map[string]interface{})
	}

	device, err := h.db.GetDevice(req.DeviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	if err := h.db.CreateCommand(req.DeviceID, req.Command, req.Parameters); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save command"})
		return
	}

	if h.mqtt == nil {
		h.logger.Warn("MQTT client not available, command not sent", zap.String("device_id", req.DeviceID), zap.String("command", req.Command))
		c.JSON(http.StatusOK, gin.H{"message": "Command saved, MQTT unavailable"})
		return
	}

	payload, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode command"})
		return
	}

	topic := fmt.Sprintf("%s/control", req.DeviceID)
	if device.Type == "smart_light" {
		topic = "smart_light/control"
	}

	if token := h.mqtt.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		h.logger.Error("Failed to send MQTT command", zap.Error(token.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send command"})
		return
	}

	h.logger.Info("Command sent", zap.String("device_id", req.DeviceID), zap.String("command", req.Command))
	metrics.RecordControlCommand(req.DeviceID, req.Command)
	metrics.RecordMQTTMessage(topic, "sent")
	c.JSON(http.StatusOK, gin.H{"message": "Command sent", "topic": topic})
}

func (h *Handler) SetThreshold(c *gin.Context) {
	var req struct {
		Threshold float64 `json:"threshold"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if h.mqtt == nil {
		h.logger.Warn("MQTT client not available, threshold not set")
		c.JSON(http.StatusOK, gin.H{"message": "Threshold saved, MQTT unavailable", "threshold": req.Threshold})
		return
	}

	topic := "smart_light/threshold"
	payload := []byte(fmt.Sprintf("%.1f", req.Threshold))

	if token := h.mqtt.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send threshold"})
		return
	}

	h.logger.Info("Threshold set", zap.Float64("threshold", req.Threshold))
	c.JSON(http.StatusOK, gin.H{"message": "Threshold set", "threshold": req.Threshold})
}

func (h *Handler) GetCommandHistory(c *gin.Context) {
	deviceID := c.Param("device_id")
	commands, err := h.db.GetCommands(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commands)
}

func (h *Handler) GetControlRecords(c *gin.Context) {
	deviceID := c.Query("device_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	if pageSize > 100 { pageSize = 100 }

	commands, total, err := h.db.GetAllCommands(deviceID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if commands == nil { commands = []*db.Command{} }

	totalPages := (total + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, gin.H{
		"records":     commands,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func (h *Handler) GetControlStats(c *gin.Context) {
	stats, err := h.db.GetCommandStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetMQTTStatus(c *gin.Context) {
	status := mqttstatus.GetStatus()

	subscriptions := []string{}
	clientActive := h.mqtt != nil
	if clientActive {
		subscriptions = []string{
			"smart_light/#",
			"+/register",
			"+/status",
			"+/control",
			"+/metric",
		}
		if h.mqtt.IsConnected() {
			status.Connected = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"connected":     status.Connected || clientActive,
		"broker":        status.Broker,
		"port":          status.Port,
		"protocol":      status.Protocol,
		"tls_enabled":   status.TLSEnabled,
		"connected_at":  status.ConnectedAt,
		"subscriptions": subscriptions,
	})
}

func (h *Handler) WebSocketHandler(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}
			if len(h.allowedOrigins) == 0 {
				// 默认策略：同源检查
				host := r.Host
				return origin == "http://"+host || origin == "https://"+host
			}
			for _, allowed := range h.allowedOrigins {
				if allowed == "*" || allowed == origin {
					return true
				}
			}
			return false
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", zap.Error(err))
		return
	}

	client := &WebSocketClient{
		hub:  h.wsHub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

type WebSocketHub struct {
	clients    map[*WebSocketClient]bool
	broadcast  chan []byte
	register   chan *WebSocketClient
	unregister chan *WebSocketClient
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		broadcast:  make(chan []byte),
		register:   make(chan *WebSocketClient),
		unregister: make(chan *WebSocketClient),
		clients:    make(map[*WebSocketClient]bool),
	}
}

func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *WebSocketHub) Send(message []byte) {
	h.broadcast <- message
}

type WebSocketClient struct {
	hub  *WebSocketHub
	conn *websocket.Conn
	send chan []byte
}

func (c *WebSocketClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *WebSocketClient) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)
		w.Close()
	}
}

type SmartLightStatus struct {
	ID        string  `json:"id"`
	Lux       float64 `json:"lux"`
	Presence  bool    `json:"presence"`
	Light     string  `json:"light"`
	Mode      string  `json:"mode"`
	Threshold int     `json:"threshold"`
	Uptime    int     `json:"uptime"`
	RSSI      int     `json:"rssi"`
	Online    bool    `json:"online"`
}

func (h *Handler) HandleMQTTMessage(msg mqtt.Message) {
	data, _ := json.Marshal(map[string]interface{}{
		"topic":   msg.Topic(),
		"payload": string(msg.Payload()),
		"time":    time.Now().Unix(),
	})
	h.wsHub.Send(data)

	topic := msg.Topic()
	metrics.RecordMQTTMessage(topic, "received")
	
	if strings.HasPrefix(topic, "smart_light/") {
		h.handleSmartLightMessage(topic, msg.Payload())
	} else {
		h.handleGenericDeviceMessage(topic, msg.Payload())
	}
}

func (h *Handler) handleSmartLightMessage(topic string, payload []byte) {
	parts := strings.Split(topic, "/")
	if len(parts) < 2 {
		return
	}

	subTopic := parts[1]

	switch subTopic {
	case "register":
		var reg struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		if json.Unmarshal(payload, &reg) == nil && reg.ID != "" {
			device, err := h.db.GetDevice(reg.ID)
			if err == nil && device != nil && device.ID != "" {
				h.db.UpdateDeviceStatus(device.ID, "online")
			} else {
				h.db.CreateDevice(&db.Device{
					ID:        reg.ID,
					Name:      reg.Name,
					Type:      "smart_light",
					Status:    "online",
					Metadata:  make(map[string]interface{}),
					CreatedAt: time.Now(),
				})
				h.db.UpdateDeviceStatus(reg.ID, "online")
			}
		}
	case "status":
		var status SmartLightStatus
		if json.Unmarshal(payload, &status) == nil {
			device, _ := h.db.GetDevice(status.ID)
			if device != nil {
				device.Status = "online"
				device.Metadata = map[string]interface{}{
					"lux":       status.Lux,
					"presence":  status.Presence,
					"light":     status.Light,
					"mode":      status.Mode,
					"threshold": status.Threshold,
					"uptime":    status.Uptime,
					"rssi":      status.RSSI,
					"online":    status.Online,
				}
				h.db.UpdateDevice(device)
			}

			h.db.CreateMetric(status.ID, "lux", status.Lux)
			h.db.CreateMetric(status.ID, "rssi", float64(status.RSSI))
			presenceVal := 0.0
			if status.Presence {
				presenceVal = 1.0
			}
			h.db.CreateMetric(status.ID, "presence", presenceVal)
		}
	case "state":
		devices, _ := h.db.GetDevices()
		for _, d := range devices {
			if d.Type == "smart_light" {
				h.db.UpdateDeviceStatus(d.ID, "online")
			}
		}
	case "presence":
		presence, _ := strconv.ParseFloat(string(payload), 64)
		devices, _ := h.db.GetDevices()
		for _, d := range devices {
			if d.Type == "smart_light" {
				if d.Metadata == nil {
					d.Metadata = make(map[string]interface{})
				}
				d.Metadata["presence"] = presence > 0
				h.db.UpdateDevice(d)
				h.db.CreateMetric(d.ID, "presence", presence)
			}
		}
	}
}

func (h *Handler) handleGenericDeviceMessage(topic string, payload []byte) {
	parts := strings.Split(topic, "/")
	if len(parts) < 2 {
		return
	}

	deviceID := parts[0]
	subTopic := parts[1]

	switch subTopic {
	case "register":
		var reg struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		}
		if json.Unmarshal(payload, &reg) == nil && reg.ID != "" {
			device, err := h.db.GetDevice(reg.ID)
			if err == nil && device != nil && device.ID != "" {
				h.db.UpdateDeviceStatus(device.ID, "online")
			} else {
				h.db.CreateDevice(&db.Device{
					ID:        reg.ID,
					Name:      reg.Name,
					Type:      reg.Type,
					Status:    "online",
					Metadata:  make(map[string]interface{}),
					CreatedAt: time.Now(),
				})
				h.db.UpdateDeviceStatus(reg.ID, "online")
			}
		}
	case "status":
		var status map[string]interface{}
		if json.Unmarshal(payload, &status) == nil {
			device, _ := h.db.GetDevice(deviceID)
			if device != nil {
				device.Status = "online"
				if status != nil {
					device.Metadata = status
				}
				h.db.UpdateDevice(device)
			}
		}
	case "metric":
		var metric struct {
			Name  string  `json:"name"`
			Value float64 `json:"value"`
		}
		if json.Unmarshal(payload, &metric) == nil {
			h.db.CreateMetric(deviceID, metric.Name, metric.Value)
		}
	}
}

// CORSMiddleware 跨域资源共享中间件
func (h *Handler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowOrigin := ""
		
		if len(h.allowedOrigins) > 0 {
			for _, allowed := range h.allowedOrigins {
				if allowed == "*" {
					allowOrigin = "*"
					break
				}
				if allowed == origin {
					allowOrigin = origin
					break
				}
			}
		} else if origin != "" {
			// 默认同源策略
			host := c.Request.Host
			if origin == "http://"+host || origin == "https://"+host {
				allowOrigin = origin
			}
		}

		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			if allowOrigin != "*" {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware 安全响应头中间件
func (h *Handler) SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}
