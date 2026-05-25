const API_BASE = '/api'

class ApiService {
  constructor() {
    this.token = null
  }

  setToken(token) {
    this.token = token
  }

  getHeaders() {
    const headers = {
      'Content-Type': 'application/json',
    }
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }
    return headers
  }

  async request(endpoint, options = {}) {
    const url = `${API_BASE}${endpoint}`
    const config = {
      headers: this.getHeaders(),
      ...options,
    }
    
    const response = await fetch(url, config)
    const data = await response.json()
    
    if (!response.ok) {
      throw new Error(data.error || '请求失败')
    }
    
    return data
  }

  // 认证相关
  async login(username, password) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
  }

  async register(username, password) {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
  }

  async updateAccount(oldUsername, newUsername, password) {
    return this.request('/auth/update-account', {
      method: 'POST',
      body: JSON.stringify({ old_username: oldUsername, new_username: newUsername, password }),
    })
  }

  async updatePassword(username, oldPassword, newPassword) {
    return this.request('/auth/update-password', {
      method: 'POST',
      body: JSON.stringify({ username, old_password: oldPassword, new_password: newPassword }),
    })
  }

  // 设备相关
  async getDevices() {
    return this.request('/devices')
  }

  async addDevice(id, name, type, note) {
    return this.request('/devices', {
      method: 'POST',
      body: JSON.stringify({
        id,
        name,
        type,
        metadata: note ? { note } : {},
      }),
    })
  }

  async deleteDevice(id) {
    return this.request(`/devices/${id}`, {
      method: 'DELETE',
    })
  }

  async getDeviceStats() {
    return this.request('/devices/stats')
  }

  async getDashboard() {
    return this.request('/dashboard')
  }

  async getMQTTStatus() {
    return this.request('/mqtt/status')
  }

  // 设备类型相关
  async getDeviceTypes() {
    return this.request('/device-types')
  }

  async addDeviceType(typeId, name, description, icon, commands) {
    return this.request('/device-types', {
      method: 'POST',
      body: JSON.stringify({
        type_id: typeId,
        name,
        description,
        icon,
        commands,
      }),
    })
  }

  // 控制相关
  async sendCommand(deviceId, command) {
    return this.request('/control/send', {
      method: 'POST',
      body: JSON.stringify({ device_id: deviceId, command }),
    })
  }

  async setThreshold(threshold) {
    return this.request('/control/threshold', {
      method: 'POST',
      body: JSON.stringify({ threshold: parseFloat(threshold) }),
    })
  }
}

export default new ApiService()
