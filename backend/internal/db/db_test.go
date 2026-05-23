package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLite_UserCRUD(t *testing.T) {
	// 创建临时数据库
	tmpFile, err := os.CreateTemp("", "test-*.db")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// 初始化数据库
	db, err := NewSQLite(tmpFile.Name())
	assert.NoError(t, err)
	defer db.Close()

	// 测试创建用户
	user, err := db.CreateUser("testuser", "passwordhash")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Positive(t, user.ID)

	// 测试通过用户名获取用户
	foundUser, err := db.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, "testuser", foundUser.Username)

	// 测试更新用户名
	err = db.UpdateUser(user.ID, "newname", "")
	assert.NoError(t, err)
	updatedUser, err := db.GetUserByUsername("newname")
	assert.NoError(t, err)
	assert.Equal(t, "newname", updatedUser.Username)

	// 测试更新密码
	err = db.UpdatePassword(user.ID, "newhash")
	assert.NoError(t, err)
	passwordUser, err := db.GetUserByUsername("newname")
	assert.NoError(t, err)
	assert.Equal(t, "newhash", passwordUser.PasswordHash)
}

func TestSQLite_DeviceCRUD(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	db, err := NewSQLite(tmpFile.Name())
	assert.NoError(t, err)
	defer db.Close()

	// 创建设备
	device := &Device{
		ID:        "device1",
		Name:      "Test Device",
		Type:      "sensor",
		Status:    "offline",
		Metadata:  map[string]interface{}{"key": "value"},
	}
	err = db.CreateDevice(device)
	assert.NoError(t, err)

	// 获取设备
	foundDevice, err := db.GetDevice("device1")
	assert.NoError(t, err)
	assert.Equal(t, "Test Device", foundDevice.Name)

	// 获取所有设备
	devices, err := db.GetDevices()
	assert.NoError(t, err)
	assert.Len(t, devices, 1)

	// 更新设备
	device.Name = "Updated Device"
	err = db.UpdateDevice(device)
	assert.NoError(t, err)
	updatedDevice, _ := db.GetDevice("device1")
	assert.Equal(t, "Updated Device", updatedDevice.Name)

	// 更新设备状态
	err = db.UpdateDeviceStatus("device1", "online")
	assert.NoError(t, err)
	statusDevice, _ := db.GetDevice("device1")
	assert.Equal(t, "online", statusDevice.Status)

	// 删除设备
	err = db.DeleteDevice("device1")
	assert.NoError(t, err)
	deletedDevices, _ := db.GetDevices()
	assert.Len(t, deletedDevices, 0)
}

func TestSQLite_DeviceTypes(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	db, err := NewSQLite(tmpFile.Name())
	assert.NoError(t, err)
	defer db.Close()

	// 添加设备类型
	err = db.AddDeviceType("custom", "Custom Device", "A custom device type", "📦", []string{"read", "write"})
	assert.NoError(t, err)

	// 获取设备类型
	types, err := db.GetDeviceTypes()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(types), 1)

	// 获取单个设备类型
	deviceType, err := db.GetDeviceType("custom")
	assert.NoError(t, err)
	assert.Equal(t, "Custom Device", deviceType["name"])
}