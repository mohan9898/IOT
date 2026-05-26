package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// 命令中英文翻译映射
var commandTranslations = map[string]string{
	"ON":     "开灯",
	"OFF":    "关灯",
	"AUTO":   "自动模式",
	"SET_THRESHOLD": "设置阈值",
	"READ":   "读取",
	"CALIBRATE": "校准",
	"START":  "启动",
	"STOP":   "停止",
	"RESET":  "重置",
	"SNAPSHOT": "拍照",
	"RECORD": "录制",
	"STREAM": "推流",
	"SET_TEMP": "设置温度",
	"SET_MODE": "设置模式",
	"TOGGLE": "切换",
}

// translateCommandToChinese 将英文命令翻译为中文
func translateCommandToChinese(command string) string {
	cmd := strings.ToUpper(strings.TrimSpace(command))
	if translation, ok := commandTranslations[cmd]; ok {
		return translation
	}
	return command // 如果没有匹配的翻译，返回原文
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Device struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	GroupID   sql.NullInt64          `json:"group_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata"`
	LastSeen  sql.NullTime           `json:"last_seen,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

type Command struct {
	ID         int64                  `json:"id"`
	DeviceID   string                 `json:"device_id"`
	Command    string                 `json:"command"`
	Parameters map[string]interface{} `json:"parameters"`
	Status     string                 `json:"status"`
	SentAt     time.Time              `json:"sent_at"`
	ExecutedAt sql.NullTime           `json:"executed_at,omitempty"`
	Response   string                 `json:"response,omitempty"`
}

type Metric struct {
	ID        int64     `json:"id"`
	DeviceID  string    `json:"device_id"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type SQLite struct {
	db *sql.DB
}

func NewSQLite(path string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	s := &SQLite{db: db}
	if err := s.initTables(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SQLite) DB() *sql.DB { return s.db }
func (s *SQLite) Close() error { return s.db.Close() }

func (s *SQLite) initTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		status TEXT DEFAULT 'offline',
		group_id INTEGER,
		metadata TEXT,
		last_seen TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS device_types (
		type_id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		icon TEXT DEFAULT '📦',
		supported_commands TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id TEXT NOT NULL,
		command TEXT NOT NULL,
		parameters TEXT,
		status TEXT DEFAULT 'pending',
		sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		executed_at TIMESTAMP,
		response TEXT,
		FOREIGN KEY (device_id) REFERENCES devices(id)
	);

	CREATE TABLE IF NOT EXISTS metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id TEXT NOT NULL,
		name TEXT NOT NULL,
		value REAL NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (device_id) REFERENCES devices(id)
	);

	INSERT OR IGNORE INTO users (username, password_hash) VALUES ('admin', ?);
	
	INSERT OR IGNORE INTO device_types (type_id, name, description, icon, supported_commands) VALUES 
	('smart_light', '智能灯', 'ESP8266智能感应灯，支持光照感应和人体检测', '💡', '["ON", "OFF", "AUTO", "SET_THRESHOLD"]'),
	('sensor', '传感器', '通用传感器设备', '📡', '["READ", "CALIBRATE"]'),
	('actuator', '执行器', '通用执行器设备', '⚙️', '["ON", "OFF", "TOGGLE"]'),
	('controller', '控制器', '智能控制器设备', '🔧', '["START", "STOP", "RESET"]'),
	('camera', '摄像头', 'IP摄像头设备', '📷', '["SNAPSHOT", "RECORD", "STREAM"]'),
	('thermostat', '恒温器', '智能温控设备', '🌡️', '["SET_TEMP", "SET_MODE", "OFF"]'),
	('switch', '智能开关', '智能开关设备', '🔌', '["ON", "OFF", "TOGGLE"]');
	`

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	_, err := s.db.Exec(schema, string(hash))
	return err
}

func (s *SQLite) GetUserByUsername(username string) (*User, error) {
	u := &User{}
	err := s.db.QueryRow("SELECT id, username, password_hash, created_at FROM users WHERE username = ?", username).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt,
	)
	return u, err
}

func (s *SQLite) CreateUser(username, passwordHash string) (*User, error) {
	res, err := s.db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &User{ID: id, Username: username, CreatedAt: time.Now()}, nil
}

func (s *SQLite) UpdateUser(id int64, username, passwordHash string) error {
	if passwordHash != "" {
		_, err := s.db.Exec("UPDATE users SET username = ?, password_hash = ? WHERE id = ?", username, passwordHash, id)
		return err
	}
	_, err := s.db.Exec("UPDATE users SET username = ? WHERE id = ?", username, id)
	return err
}

func (s *SQLite) UpdatePassword(id int64, newPasswordHash string) error {
	_, err := s.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", newPasswordHash, id)
	return err
}

func (s *SQLite) GetDevices() ([]*Device, error) {
	rows, err := s.db.Query("SELECT id, name, type, status, group_id, metadata, last_seen, created_at FROM devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*Device
	for rows.Next() {
		d := &Device{}
		var metaStr string
		if err := rows.Scan(&d.ID, &d.Name, &d.Type, &d.Status, &d.GroupID, &metaStr, &d.LastSeen, &d.CreatedAt); err != nil {
			continue
		}
		json.Unmarshal([]byte(metaStr), &d.Metadata)
		devices = append(devices, d)
	}
	return devices, nil
}

func (s *SQLite) GetDevice(id string) (*Device, error) {
	d := &Device{}
	var metaStr string
	err := s.db.QueryRow("SELECT id, name, type, status, group_id, metadata, last_seen, created_at FROM devices WHERE id = ?", id).Scan(
		&d.ID, &d.Name, &d.Type, &d.Status, &d.GroupID, &metaStr, &d.LastSeen, &d.CreatedAt,
	)
	if err == nil {
		json.Unmarshal([]byte(metaStr), &d.Metadata)
	}
	return d, err
}

func (s *SQLite) CreateDevice(d *Device) error {
	meta, _ := json.Marshal(d.Metadata)
	_, err := s.db.Exec(`INSERT OR IGNORE INTO devices (id, name, type, status, metadata, last_seen) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		d.ID, d.Name, d.Type, d.Status, string(meta))
	return err
}

func (s *SQLite) UpdateDevice(d *Device) error {
	meta, _ := json.Marshal(d.Metadata)
	_, err := s.db.Exec(`UPDATE devices SET name = ?, type = ?, status = ?, metadata = ?, last_seen = CURRENT_TIMESTAMP WHERE id = ?`,
		d.Name, d.Type, d.Status, string(meta), d.ID)
	return err
}

func (s *SQLite) UpdateDeviceStatus(id, status string) error {
	_, err := s.db.Exec("UPDATE devices SET status = ?, last_seen = CURRENT_TIMESTAMP WHERE id = ?", status, id)
	return err
}

func (s *SQLite) MarkOfflineDevices(timeoutMinutes int) (int, error) {
	result, err := s.db.Exec(
		"UPDATE devices SET status = 'offline' WHERE status = 'online' AND last_seen < datetime('now', ?)",
		fmt.Sprintf("-%d minutes", timeoutMinutes),
	)
	if err != nil {
		return 0, err
	}
	count, _ := result.RowsAffected()
	return int(count), nil
}

func (s *SQLite) UpdateDeviceMetadata(id string, metadata map[string]interface{}) error {
	meta, _ := json.Marshal(metadata)
	_, err := s.db.Exec("UPDATE devices SET metadata = ? WHERE id = ?", string(meta), id)
	return err
}

func (s *SQLite) DeleteDevice(id string) error {
	_, err := s.db.Exec("DELETE FROM devices WHERE id = ?", id)
	return err
}

func (s *SQLite) GetDeviceStats() (total, online, offline int, err error) {
	s.db.QueryRow("SELECT COUNT(*) FROM devices").Scan(&total)
	s.db.QueryRow("SELECT COUNT(*) FROM devices WHERE status = 'online'").Scan(&online)
	s.db.QueryRow("SELECT COUNT(*) FROM devices WHERE status = 'offline'").Scan(&offline)
	return total, online, offline, nil
}

func (s *SQLite) GetDeviceTypes() ([]map[string]interface{}, error) {
	rows, err := s.db.Query("SELECT type_id, name, description, icon, supported_commands FROM device_types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []map[string]interface{}
	for rows.Next() {
		var typeID, name, desc, icon, cmds string
		rows.Scan(&typeID, &name, &desc, &icon, &cmds)
		types = append(types, map[string]interface{}{
			"type_id":             typeID,
			"name":               name,
			"description":        desc,
			"icon":               icon,
			"supported_commands": cmds,
		})
	}
	return types, nil
}

func (s *SQLite) AddDeviceType(typeID, name, desc, icon string, commands []string) error {
	cmdsJSON, _ := json.Marshal(commands)
	_, err := s.db.Exec(`INSERT OR REPLACE INTO device_types 
		(type_id, name, description, icon, supported_commands) VALUES (?, ?, ?, ?, ?)`,
		typeID, name, desc, icon, string(cmdsJSON))
	return err
}

func (s *SQLite) GetDeviceType(typeID string) (map[string]interface{}, error) {
	var name, desc, icon, cmds string
	err := s.db.QueryRow("SELECT name, description, icon, supported_commands FROM device_types WHERE type_id = ?", typeID).Scan(
		&name, &desc, &icon, &cmds)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"type_id":             typeID,
		"name":               name,
		"description":        desc,
		"icon":               icon,
		"supported_commands": cmds,
	}, nil
}

func (s *SQLite) CreateCommand(deviceID, command string, params map[string]interface{}) error {
	p, _ := json.Marshal(params)
	_, err := s.db.Exec("INSERT INTO commands (device_id, command, parameters) VALUES (?, ?, ?)", deviceID, command, string(p))
	if err != nil {
		return err
	}
	s.CleanupOldCommands(200)
	return nil
}

func (s *SQLite) CleanupOldCommands(maxRows int) {
	s.db.Exec(`
		DELETE FROM commands WHERE id NOT IN (
			SELECT id FROM commands ORDER BY sent_at DESC LIMIT ?
		)
	`, maxRows)
}

func (s *SQLite) GetCommands(deviceID string) ([]*Command, error) {
	rows, err := s.db.Query("SELECT id, device_id, command, parameters, status, sent_at, executed_at, response FROM commands WHERE device_id = ? ORDER BY sent_at DESC", deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cmds []*Command
	for rows.Next() {
		c := &Command{}
		var paramsStr string
		var originalCommand string
		rows.Scan(&c.ID, &c.DeviceID, &originalCommand, &paramsStr, &c.Status, &c.SentAt, &c.ExecutedAt, &c.Response)
		json.Unmarshal([]byte(paramsStr), &c.Parameters)
		// 将命令翻译为中文显示
		c.Command = translateCommandToChinese(originalCommand)
		cmds = append(cmds, c)
	}
	return cmds, nil
}

func (s *SQLite) UpdateCommandStatus(id int64, status, response string) error {
	_, err := s.db.Exec("UPDATE commands SET status = ?, response = ?, executed_at = CURRENT_TIMESTAMP WHERE id = ?", status, response, id)
	return err
}

func (s *SQLite) GetAllCommands(deviceID string, page, pageSize int) ([]*Command, int, error) {
	where := ""
	args := []interface{}{}
	if deviceID != "" {
		where = "WHERE device_id = ?"
		args = append(args, deviceID)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM commands " + where
	s.db.QueryRow(countQuery, args...).Scan(&total)

	offset := (page - 1) * pageSize
	query := "SELECT id, device_id, command, parameters, status, sent_at, executed_at, response FROM commands " + where + " ORDER BY sent_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cmds []*Command
	for rows.Next() {
		c := &Command{}
		var paramsStr string
		var originalCommand string
		rows.Scan(&c.ID, &c.DeviceID, &originalCommand, &paramsStr, &c.Status, &c.SentAt, &c.ExecutedAt, &c.Response)
		json.Unmarshal([]byte(paramsStr), &c.Parameters)
		// 将命令翻译为中文显示
		c.Command = translateCommandToChinese(originalCommand)
		cmds = append(cmds, c)
	}
	return cmds, total, nil
}

func (s *SQLite) GetCommandStats() (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	var totalCount int
	s.db.QueryRow("SELECT COUNT(*) FROM commands").Scan(&totalCount)
	stats["total"] = totalCount

	var todayCount int
	s.db.QueryRow("SELECT COUNT(*) FROM commands WHERE date(sent_at) = date('now')").Scan(&todayCount)
	stats["today"] = todayCount

	rows, err := s.db.Query("SELECT device_id, COUNT(*) as cnt FROM commands GROUP BY device_id ORDER BY cnt DESC LIMIT 10")
	if err == nil {
		defer rows.Close()
		byDevice := []map[string]interface{}{}
		for rows.Next() {
			var did string
			var cnt int
			rows.Scan(&did, &cnt)
			byDevice = append(byDevice, map[string]interface{}{"device_id": did, "count": cnt})
		}
		stats["by_device"] = byDevice
	}

	var successCount int
	s.db.QueryRow("SELECT COUNT(*) FROM commands WHERE status = 'success'").Scan(&successCount)
	stats["success"] = successCount

	var failedCount int
	s.db.QueryRow("SELECT COUNT(*) FROM commands WHERE status = 'failed'").Scan(&failedCount)
	stats["failed"] = failedCount

	return stats, nil
}

func (s *SQLite) CreateMetric(deviceID, name string, value float64) error {
	_, err := s.db.Exec("INSERT INTO metrics (device_id, name, value) VALUES (?, ?, ?)", deviceID, name, value)
	return err
}

func (s *SQLite) GetMetrics(deviceID, name string, limit int) ([]*Metric, error) {
	rows, err := s.db.Query("SELECT id, device_id, name, value, timestamp FROM metrics WHERE device_id = ? AND name = ? ORDER BY timestamp DESC LIMIT ?", deviceID, name, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*Metric
	for rows.Next() {
		m := &Metric{}
		rows.Scan(&m.ID, &m.DeviceID, &m.Name, &m.Value, &m.Timestamp)
		metrics = append(metrics, m)
	}
	return metrics, nil
}
