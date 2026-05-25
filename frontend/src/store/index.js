import { reactive } from 'vue'
import api from '../services/api'

const store = reactive({
  // 暴露 api 服务给子组件
  api,

  // 状态
  loggedIn: false,
  token: '',
  currentUsername: '',
  devices: [],
  deviceTypes: [],
  stats: { total: 0, online: 0, offline: 0 },
  selectedDevice: null,
  ws: null,
  notificationsEnabled: false,
  seenOnboarding: false,

  // 加载状态
  loading: {
    devices: false,
    deviceTypes: false,
    stats: false,
  },

  // 方法
  setToken(token) {
    this.token = token
    api.setToken(token)
  },

  async login(username, password) {
    const data = await api.login(username, password)
    this.token = data.token
    this.currentUsername = data.user.username
    this.loggedIn = true
    api.setToken(data.token)
    localStorage.setItem('iot_token', data.token)
    localStorage.setItem('iot_username', data.user.username)
    
    await Promise.all([
      this.loadDevices(),
      this.loadStats(),
      this.loadDeviceTypes(),
    ])
    
    this.connectWebSocket()
    return data
  },

  async register(username, password) {
    const data = await api.register(username, password)
    this.token = data.token
    this.currentUsername = data.user.username
    this.loggedIn = true
    api.setToken(data.token)
    localStorage.setItem('iot_token', data.token)
    localStorage.setItem('iot_username', data.user.username)
    
    await Promise.all([
      this.loadDevices(),
      this.loadStats(),
      this.loadDeviceTypes(),
    ])
    
    this.connectWebSocket()
    return data
  },

  logout() {
    this.loggedIn = false
    this.token = ''
    this.currentUsername = ''
    this.devices = []
    this.stats = { total: 0, online: 0, offline: 0 }
    this.selectedDevice = null
    localStorage.removeItem('iot_token')
    localStorage.removeItem('iot_username')
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    api.setToken(null)
  },

  async tryAutoLogin() {
    const savedToken = localStorage.getItem('iot_token')
    const savedUsername = localStorage.getItem('iot_username')
    if (!savedToken) return false
    
    this.token = savedToken
    this.currentUsername = savedUsername || ''
    api.setToken(savedToken)
    
    try {
      await this.loadDevices()
      this.loggedIn = true
      this.connectWebSocket()
      return true
    } catch (e) {
      this.logout()
      return false
    }
  },

  async loadDevices() {
    this.loading.devices = true
    try {
      const data = await api.getDevices()
      this.devices = data || []
    } finally {
      this.loading.devices = false
    }
  },

  async loadDeviceTypes() {
    this.loading.deviceTypes = true
    try {
      if (this.loggedIn) {
        const data = await api.getDeviceTypes()
        this.deviceTypes = data || []
      }
    } finally {
      this.loading.deviceTypes = false
    }
  },

  async loadStats() {
    this.loading.stats = true
    try {
      const data = await api.getDeviceStats()
      this.stats = data || { total: 0, online: 0, offline: 0 }
    } finally {
      this.loading.stats = false
    }
  },

  selectDevice(device) {
    this.selectedDevice = device
  },

  closeDeviceModal() {
    this.selectedDevice = null
  },

  getDeviceIcon(type) {
    const dt = this.deviceTypes.find(t => t.type_id === type)
    return dt?.icon || '📦'
  },

  getDeviceTypeName(type) {
    const dt = this.deviceTypes.find(t => t.type_id === type)
    return dt?.name || type
  },

  formatUptime(seconds) {
    const hours = Math.floor(seconds / 3600)
    const mins = Math.floor((seconds % 3600) / 60)
    return `${hours}h ${mins}m`
  },

  connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    this.ws = new WebSocket(`${protocol}//${host}/api/ws?token=${this.token}`)
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      if (data.topic && data.topic.includes('status')) {
        const oldDevices = [...this.devices]
        this.loadDevices()
        this.loadStats()
        
        if (typeof data.payload === 'string') {
          try {
            const payload = JSON.parse(data.payload)
            this.checkDeviceAlerts(oldDevices, payload)
          } catch (e) {}
        }
      }
    }
    this.ws.onclose = () => {
      setTimeout(() => {
        if (this.loggedIn) {
          this.connectWebSocket()
        }
      }, 5000)
    }
  },

  checkDeviceAlerts(oldDevices, mqttPayload) {
    if (!this.notificationsEnabled) return
    if (!mqttPayload.id) return

    const oldDevice = oldDevices.find(d => d.id === mqttPayload.id)
    if (!oldDevice) return

    const wasOnline = oldDevice.status === 'online'
    const isOnline = mqttPayload.online !== false && mqttPayload.status !== 'offline'

    if (wasOnline && !isOnline) {
      this.sendNotification(
        '设备离线告警',
        `设备 "${mqttPayload.id}" 已离线`,
        '⚠️'
      )
    } else if (!wasOnline && isOnline && oldDevice.status !== undefined) {
      this.sendNotification(
        '设备上线',
        `设备 "${mqttPayload.id}" 已恢复在线`,
        '✅'
      )
    }
  },

  sendNotification(title, body, icon) {
    if (typeof Notification === 'undefined') return
    
    if (Notification.permission === 'granted') {
      new Notification(title, { body, icon: 'data:image/svg+xml,' + encodeURIComponent('<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32"><rect width="32" height="32" rx="6" fill="#3B82F6"/><text x="16" y="22" font-size="20" text-anchor="middle">' + icon + '</text></svg>'), tag: 'iot-device' })
    } else if (Notification.permission !== 'denied') {
      Notification.requestPermission()
    }
  },

  enableNotifications() {
    this.notificationsEnabled = true
    if (typeof Notification !== 'undefined' && Notification.permission !== 'granted') {
      Notification.requestPermission()
    }
  },

  disableNotifications() {
    this.notificationsEnabled = false
  },
})

export default store
