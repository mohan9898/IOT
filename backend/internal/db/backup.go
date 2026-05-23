package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

// BackupManager 处理数据库备份
type BackupManager struct {
	db           *SQLite
	dbPath       string
	backupPath   string
	backupHours  int
	logger       *zap.Logger
	stopChan     chan struct{}
	backupChan   chan struct{}
}

// NewBackupManager 创建备份管理器
func NewBackupManager(db *SQLite, dbPath, backupPath string, backupHours int, logger *zap.Logger) *BackupManager {
	return &BackupManager{
		db:          db,
		dbPath:      dbPath,
		backupPath:  backupPath,
		backupHours: backupHours,
		logger:      logger,
		stopChan:    make(chan struct{}),
		backupChan:  make(chan struct{}),
	}
}

// Start 启动定期备份
func (b *BackupManager) Start() {
	// 确保备份目录存在
	if err := os.MkdirAll(b.backupPath, 0700); err != nil {
		b.logger.Error("Failed to create backup directory", zap.Error(err))
		return
	}

	go b.run()
}

// Stop 停止备份服务
func (b *BackupManager) Stop() {
	close(b.stopChan)
}

// TriggerBackup 手动触发备份
func (b *BackupManager) TriggerBackup() {
	select {
	case b.backupChan <- struct{}{}:
	default:
		b.logger.Debug("Backup already queued")
	}
}

func (b *BackupManager) run() {
	b.logger.Info("Backup manager started", zap.Int("interval_hours", b.backupHours))

	// 立即执行一次备份
	b.doBackup()

	ticker := time.NewTicker(time.Duration(b.backupHours) * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.doBackup()
		case <-b.backupChan:
			b.doBackup()
		case <-b.stopChan:
			b.logger.Info("Backup manager stopped")
			return
		}
	}
}

func (b *BackupManager) doBackup() {
	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(b.backupPath, fmt.Sprintf("iot-%s.db", timestamp))

	// 使用 SQLite 的 VACUUM INTO 命令进行安全备份
	_, err := b.db.DB().Exec(fmt.Sprintf("VACUUM INTO '%s'", backupFile))
	if err != nil {
		b.logger.Error("Failed to backup database using VACUUM", zap.Error(err))
		return
	}

	b.logger.Info("Backup completed successfully", zap.String("file", backupFile))

	// 清理旧备份（保留最近 7 天）
	b.cleanOldBackups()
}

func (b *BackupManager) cleanOldBackups() {
	entries, err := os.ReadDir(b.backupPath)
	if err != nil {
		b.logger.Error("Failed to read backup directory", zap.Error(err))
		return
	}

	cutoff := time.Now().Add(-7 * 24 * time.Hour)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			oldFile := filepath.Join(b.backupPath, entry.Name())
			if err := os.Remove(oldFile); err != nil {
				b.logger.Error("Failed to remove old backup", zap.String("file", oldFile), zap.Error(err))
			} else {
				b.logger.Info("Old backup removed", zap.String("file", oldFile))
			}
		}
	}
}

// ListBackups 列出可用的备份
func (b *BackupManager) ListBackups() ([]string, error) {
	var backups []string

	entries, err := os.ReadDir(b.backupPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}

// Restore 从备份恢复数据库
func (b *BackupManager) Restore(backupName string) error {
	// 防止路径遍历攻击：检查文件名是否包含路径分隔符
	if strings.Contains(backupName, "/") || strings.Contains(backupName, "\\") || strings.Contains(backupName, "..") {
		return fmt.Errorf("invalid backup name")
	}
	backupFile := filepath.Join(b.backupPath, backupName)

	// 检查备份文件是否存在
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	// 备份当前数据库
	if _, err := os.Stat(b.dbPath); err == nil {
		timestamp := time.Now().Format("20060102-150405")
		oldBackup := filepath.Join(b.backupPath, fmt.Sprintf("iot-pre-restore-%s.db", timestamp))
		_, err := b.db.DB().Exec(fmt.Sprintf("VACUUM INTO '%s'", oldBackup))
		if err != nil {
			b.logger.Error("Failed to backup current database before restore", zap.Error(err))
		}
	}

	// 关闭当前数据库连接
	if err := b.db.Close(); err != nil {
		b.logger.Warn("Failed to close database connection before restore", zap.Error(err))
	}

	// 替换数据库文件
	if err := os.Rename(backupFile, b.dbPath); err != nil {
		// 尝试复制文件
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// 重新打开数据库连接
	newDb, err := NewSQLite(b.dbPath)
	if err != nil {
		return fmt.Errorf("failed to re-open database after restore: %w", err)
	}

	// 更新备份管理器的数据库引用
	*b.db = *newDb

	b.logger.Info("Database restored successfully", zap.String("backup", backupName))
	return nil
}