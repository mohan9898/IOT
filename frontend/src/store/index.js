import { reactive } from 'vue'
import api from '../services/api'

const store = reactive({
  // 状态
  loggedIn: false,
  token: '',
  currentUsername: '',
  devices: [],
  deviceTypes: [],
  stats: { total: 0, online: 0, offline: 0 },
  selectedDevice: null,
  ws: null,
  
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
    
    // 加载初始数据
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
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    api.setToken(null)
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
        this.loadDevices()
        this.loadStats()
      }
    }
  },
})

export default store
