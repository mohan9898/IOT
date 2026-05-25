<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-xl sm:text-2xl font-bold text-gray-800 whitespace-nowrap">仪表盘</h1>
      <button
        @click="refresh"
        :disabled="loading"
        class="px-4 py-2 text-sm bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-colors disabled:opacity-50"
      >
        {{ loading ? '刷新中...' : '刷新数据' }}
      </button>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs text-gray-500 mb-0.5 whitespace-nowrap">设备总数</p>
            <p class="text-3xl font-bold text-gray-800">{{ dashboard.stats.total }}</p>
          </div>
          <div class="w-10 h-10 bg-blue-100 rounded-xl flex items-center justify-center text-xl flex-shrink-0">
            📱
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs text-gray-500 mb-0.5 whitespace-nowrap">在线设备</p>
            <p class="text-3xl font-bold text-green-600">{{ dashboard.stats.online }}</p>
          </div>
          <div class="w-10 h-10 bg-green-100 rounded-xl flex items-center justify-center text-xl flex-shrink-0">
            🟢
          </div>
        </div>
        <div class="mt-2 w-full bg-gray-100 rounded-full h-1.5">
          <div
            class="bg-green-500 h-1.5 rounded-full transition-all duration-500"
            :style="{ width: percentage(dashboard.stats.online, dashboard.stats.total) + '%' }"
          ></div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs text-gray-500 mb-0.5 whitespace-nowrap">离线设备</p>
            <p class="text-3xl font-bold text-red-500">{{ dashboard.stats.offline }}</p>
          </div>
          <div class="w-10 h-10 bg-red-100 rounded-xl flex items-center justify-center text-xl flex-shrink-0">
            🔴
          </div>
        </div>
        <div class="mt-2 w-full bg-gray-100 rounded-full h-1.5">
          <div
            class="bg-red-500 h-1.5 rounded-full transition-all duration-500"
            :style="{ width: percentage(dashboard.stats.offline, dashboard.stats.total) + '%' }"
          ></div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
      <div class="flex items-center justify-between mb-2">
        <h2 class="text-lg font-semibold text-gray-800">MQTT 连接状态</h2>
        <button
          @click="refreshMQTT"
          :disabled="mqttLoading"
          class="text-sm text-blue-500 hover:text-blue-600 disabled:opacity-50"
        >
          {{ mqttLoading ? '检测中...' : '刷新检测' }}
        </button>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
        <div class="flex items-center gap-2 p-2.5 bg-gray-50 rounded-xl">
          <span
            :class="['w-2.5 h-2.5 rounded-full flex-shrink-0', mqttStatus.connected ? 'bg-green-500 animate-pulse' : 'bg-red-500']"
          ></span>
          <div class="min-w-0">
            <p class="text-xs text-gray-400 whitespace-nowrap">连接状态</p>
            <p :class="['text-xs font-semibold whitespace-nowrap', mqttStatus.connected ? 'text-green-600' : 'text-red-500']">
              {{ mqttStatus.connected ? '已连接' : '已断开' }}
            </p>
          </div>
        </div>
        <div class="p-2.5 bg-gray-50 rounded-xl min-w-0">
          <p class="text-xs text-gray-400 whitespace-nowrap">Broker 地址</p>
          <p class="text-xs font-semibold text-gray-700 truncate" :title="mqttStatus.broker">{{ mqttStatus.broker || '-' }}</p>
        </div>
        <div class="p-2.5 bg-gray-50 rounded-xl">
          <p class="text-xs text-gray-400 whitespace-nowrap">连接类型</p>
          <p class="text-xs font-semibold text-gray-700 whitespace-nowrap">
            {{ mqttStatus.connection_type || '-' }}
          </p>
        </div>
        <div class="p-2.5 bg-gray-50 rounded-xl">
          <p class="text-xs text-gray-400 whitespace-nowrap">订阅主题</p>
          <p class="text-xs font-semibold text-gray-700 whitespace-nowrap">
            {{ (mqttStatus.subscriptions || []).length }} 个
          </p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
      <div class="lg:col-span-3 bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">设备类型分布</h2>
        <div v-if="dashboard.typeDistribution.length === 0" class="text-center py-12 text-gray-400">
          <p class="text-5xl mb-3">📊</p>
          <p>暂无设备数据</p>
          <p class="text-sm mt-1">添加设备后将在此显示统计图表</p>
        </div>
        <div v-else class="space-y-4">
          <div
            v-for="item in sortedTypes"
            :key="item.type_id"
            class="flex items-center gap-3"
          >
            <span class="text-xl w-8 flex-shrink-0">{{ item.icon }}</span>
            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-center mb-1">
                <span class="text-sm font-medium text-gray-700 truncate">{{ item.name }}</span>
                <span class="text-sm text-gray-500 ml-2 flex-shrink-0">{{ item.count }} 台</span>
              </div>
              <div class="w-full bg-gray-100 rounded-full h-2">
                <div
                  class="bg-blue-500 h-2 rounded-full transition-all duration-700"
                  :style="{ width: percentage(item.count, maxCount) + '%' }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="lg:col-span-2 bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">最近设备</h2>
        <div v-if="dashboard.recentDevices.length === 0" class="text-center py-12 text-gray-400">
          <p class="text-5xl mb-3">📱</p>
          <p>暂无设备</p>
          <p class="text-sm mt-1">点击"添加设备"开始管理</p>
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="device in dashboard.recentDevices"
            :key="device.id"
            class="flex items-center justify-between p-3 bg-gray-50 rounded-xl hover:bg-gray-100 transition-colors"
          >
            <div class="flex items-center gap-3 min-w-0">
              <span class="text-lg flex-shrink-0">{{ getIcon(device.type) }}</span>
              <div class="min-w-0">
                <p class="text-sm font-medium text-gray-800 truncate">{{ device.name }}</p>
                <p class="text-xs text-gray-400 truncate">{{ device.id }}</p>
              </div>
            </div>
            <span
              :class="[
                'px-2.5 py-1 rounded-full text-xs font-medium flex-shrink-0 ml-2',
                device.status === 'online' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
              ]"
            >
              {{ device.status === 'online' ? '在线' : '离线' }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
      <h2 class="text-lg font-semibold text-gray-800 mb-4">快捷操作</h2>
      <div class="flex flex-wrap gap-3">
        <button
          @click="$emit('navigate', 'add')"
          class="px-5 py-3 bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-colors font-medium flex items-center gap-2"
        >
          <span>➕</span> 添加设备
        </button>
        <button
          @click="$emit('navigate', 'devices')"
          class="px-5 py-3 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition-colors font-medium flex items-center gap-2"
        >
          <span>📱</span> 查看所有设备
        </button>
        <button
          @click="$emit('navigate', 'types')"
          class="px-5 py-3 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition-colors font-medium flex items-center gap-2"
        >
          <span>📦</span> 管理设备类型
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted, ref } from 'vue'
import store from '../store'

const emit = defineEmits(['navigate'])

const loading = ref(false)
const mqttLoading = ref(false)

const dashboard = reactive({
  stats: { total: 0, online: 0, offline: 0 },
  typeDistribution: [],
  recentDevices: [],
})

const mqttStatus = reactive({
  connected: false,
  broker: '',
  port: 0,
  protocol: '',
  tls_enabled: false,
  connection_type: '',
  subscriptions: [],
})

const sortedTypes = computed(() => {
  return [...dashboard.typeDistribution].sort((a, b) => b.count - a.count)
})

const maxCount = computed(() => {
  if (dashboard.typeDistribution.length === 0) return 1
  return Math.max(...dashboard.typeDistribution.map(t => t.count), 1)
})

function percentage(value, total) {
  if (total === 0) return 0
  return Math.round((value / total) * 100)
}

function getIcon(type) {
  const dt = store.deviceTypes.find(t => t.type_id === type)
  return dt?.icon || '📦'
}

async function loadDashboard() {
  loading.value = true
  try {
    const data = await store.api.getDashboard()
    if (data) {
      dashboard.stats = data.stats
      dashboard.typeDistribution = data.type_distribution || []
      dashboard.recentDevices = data.recent_devices || []
    }
  } catch (e) {
    console.error('加载仪表盘失败:', e)
  } finally {
    loading.value = false
  }
}

function refresh() {
  loadDashboard()
  loadMQTTStatus()
}

async function loadMQTTStatus() {
  mqttLoading.value = true
  try {
    const data = await store.api.getMQTTStatus()
    if (data) {
      Object.assign(mqttStatus, data)
    }
  } catch (e) {
    console.error('加载MQTT状态失败:', e)
  } finally {
    mqttLoading.value = false
  }
}

function refreshMQTT() {
  loadMQTTStatus()
}

onMounted(() => {
  loadDashboard()
  loadMQTTStatus()
})

defineExpose({ refresh })
</script>