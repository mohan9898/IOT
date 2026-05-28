<template>
  <div v-if="device" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" @click.self="close">
    <div class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
      <!-- 头部 -->
      <div class="p-5 border-b border-gray-100">
        <div class="flex justify-between items-start">
          <div class="flex items-center min-w-0">
            <span class="text-4xl mr-3 flex-shrink-0">{{ icon }}</span>
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <h2 class="text-xl font-bold text-gray-800 truncate">{{ device.name }}</h2>
                <span
                  :class="['px-2 py-0.5 rounded-full text-xs font-semibold whitespace-nowrap flex-shrink-0',
                    device.status === 'online' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700']"
                >
                  {{ device.status === 'online' ? '在线' : '离线' }}
                </span>
              </div>
              <p class="text-xs text-gray-500 font-mono truncate">{{ device.id }}</p>
            </div>
          </div>
          <button @click="close" class="text-gray-400 hover:text-gray-600 text-2xl p-2 hover:bg-gray-100 rounded-full transition-colors">
            ✕
          </button>
        </div>
      </div>

      <!-- 内容 -->
      <div :class="['p-6 space-y-6 transition-all duration-300', device.status === 'offline' ? 'opacity-70' : '']">
        <!-- 智能灯特殊界面 -->
        <div v-if="device.type === 'smart_light'" class="space-y-6">
          <!-- 当前工作状态 -->
          <div class="grid grid-cols-2 gap-4">
            <div :class="['rounded-xl p-4 text-center border-2 transition-all', lightOn ? 'bg-green-50 border-green-400' : 'bg-gray-100 border-gray-300']">
              <div class="text-xs font-medium mb-1 whitespace-nowrap" :class="lightOn ? 'text-green-600' : 'text-gray-500'">灯光</div>
              <div class="text-2xl">{{ lightOn ? '🔆' : '🌙' }}</div>
              <div class="text-base font-bold mt-1 whitespace-nowrap" :class="lightOn ? 'text-green-700' : 'text-gray-400'">{{ lightOn ? '已开启' : '已关闭' }}</div>
            </div>
            <div :class="['rounded-xl p-4 text-center border-2 transition-all', isAutoMode ? 'bg-blue-50 border-blue-400' : 'bg-orange-50 border-orange-400']">
              <div class="text-xs font-medium mb-1 whitespace-nowrap" :class="isAutoMode ? 'text-blue-600' : 'text-orange-600'">模式</div>
              <div class="text-2xl">{{ isAutoMode ? '🤖' : '✋' }}</div>
              <div class="text-base font-bold mt-1 whitespace-nowrap" :class="isAutoMode ? 'text-blue-700' : 'text-orange-700'">{{ isAutoMode ? '自动控制' : '手动控制' }}</div>
            </div>
          </div>
          <!-- 传感器数据 -->
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-gradient-to-br from-yellow-50 to-yellow-100 rounded-xl p-4">
              <div class="text-xs text-yellow-600 mb-1 font-medium whitespace-nowrap">光照强度</div>
              <div class="text-2xl font-bold text-yellow-700">{{ device.metadata?.lux || '0' }} <span class="text-xs font-normal">lux</span></div>
            </div>
            <div class="bg-gradient-to-br from-purple-50 to-purple-100 rounded-xl p-4">
              <div class="text-xs text-purple-600 mb-1 font-medium whitespace-nowrap">人体感应</div>
              <div :class="['text-2xl font-bold', device.metadata?.presence ? 'text-purple-700' : 'text-gray-400']">
                {{ device.metadata?.presence ? '检测到' : '无' }}
              </div>
            </div>
            <div class="bg-gradient-to-br from-orange-50 to-orange-100 rounded-xl p-4">
              <div class="text-xs text-orange-600 mb-1 font-medium whitespace-nowrap">运行时长</div>
              <div class="text-2xl font-bold text-orange-700">{{ formatUptime(device.metadata?.uptime || 0) }}</div>
            </div>
            <div class="bg-gradient-to-br from-cyan-50 to-cyan-100 rounded-xl p-4">
              <div class="text-xs text-cyan-600 mb-1 font-medium whitespace-nowrap">WiFi信号</div>
              <div class="text-2xl font-bold text-cyan-700">{{ device.metadata?.rssi || '--' }} <span class="text-xs font-normal">dBm</span></div>
            </div>
          </div>
          
          <!-- 设备信息 -->
          <div class="bg-gray-50 rounded-xl p-4">
            <h4 class="font-semibold text-gray-700 mb-3">设备详情</h4>
            <div class="grid grid-cols-2 gap-3">
              <div v-if="device.metadata?.version" class="bg-white rounded-lg p-3">
                <div class="text-xs text-gray-500 uppercase font-medium">固件版本</div>
                <div class="text-sm font-semibold text-gray-800">{{ device.metadata.version }}</div>
              </div>
              <div v-if="device.metadata?.free_heap" class="bg-white rounded-lg p-3">
                <div class="text-xs text-gray-500 uppercase font-medium">可用内存</div>
                <div class="text-sm font-semibold text-gray-800">{{ device.metadata.free_heap }} bytes</div>
              </div>
            </div>
          </div>

          <!-- 快捷控制 -->
          <div class="grid grid-cols-3 gap-3">
            <button
              @click="sendCommand(device.type === 'pc_controller' ? 'POWER' : 'ON')"
              :disabled="loading || device.status === 'offline'"
              :class="['py-4 rounded-xl transition-all font-semibold border-2', 
                (loading || device.status === 'offline') ? 'opacity-50 cursor-not-allowed' : '',
                lightOn ? 'bg-green-500 text-white border-green-500' : 'bg-white text-green-600 border-green-300 hover:bg-green-50']"
            >
              {{ device.type === 'pc_controller' ? '🔌 开机' : '🔆 开灯' }}
            </button>
            <button
              @click="sendCommand('OFF')"
              :disabled="loading || device.status === 'offline'"
              :class="['py-4 rounded-xl transition-all font-semibold border-2',
                (loading || device.status === 'offline') ? 'opacity-50 cursor-not-allowed' : '',
                !lightOn ? 'bg-gray-500 text-white border-gray-500' : 'bg-white text-gray-600 border-gray-300 hover:bg-gray-50']"
            >
              {{ device.type === 'pc_controller' ? '🔌 断电' : '🌙 关灯' }}
            </button>
            <button
              v-if="device.type === 'smart_light'"
              @click="sendCommand('AUTO')"
              :disabled="loading || device.status === 'offline'"
              :class="['py-4 rounded-xl transition-all font-semibold border-2',
                (loading || device.status === 'offline') ? 'opacity-50 cursor-not-allowed' : '',
                isAutoMode ? 'bg-blue-500 text-white border-blue-500' : 'bg-white text-blue-600 border-blue-300 hover:bg-blue-50']"
            >
              🤖 自动
            </button>
            <button
              v-else-if="device.type === 'pc_controller'"
              @click="sendCommand('RESET')"
              :disabled="loading || device.status === 'offline'"
              :class="['py-4 rounded-xl transition-all font-semibold border-2 bg-orange-100 text-orange-700 border-orange-300 hover:bg-orange-50',
                (loading || device.status === 'offline') ? 'opacity-50 cursor-not-allowed' : '']"
            >
              🔄 重启PC
            </button>
          </div>

          <!-- 高级控制 -->
          <div class="grid grid-cols-3 gap-3 mt-3">
            <button
              @click="sendCommand('RESTART')"
              :disabled="loading || device.status === 'offline'"
              :class="['py-3 rounded-xl transition-all font-medium bg-yellow-100 text-yellow-700 border-2 border-yellow-300 hover:bg-yellow-50',
                (loading || device.status === 'offline') ? 'opacity-50 cursor-not-allowed' : '']"
            >
              🔄 远程重启
            </button>
            <button
              @click="fetchLogs"
              :disabled="loading || device.status === 'offline'"
              :class="['py-3 rounded-xl transition-all font-medium bg-purple-100 text-purple-700 border-2 border-purple-300 hover:bg-purple-50',
                (loading || device.status === 'offline') ? 'opacity-50 cursor-not-allowed' : '']"
            >
              📋 获取日志
            </button>
            <button
              @click="showOtaModal = true"
              :disabled="loading || device.status === 'offline'"
              :class="['py-3 rounded-xl transition-all font-medium bg-cyan-100 text-cyan-700 border-2 border-cyan-300 hover:bg-cyan-50',
                (loading || device.status === 'offline') ? 'opacity-50 cursor-not-allowed' : '']"
            >
              ⬆️ OTA升级
            </button>
          </div>

          <!-- 光照阈值设置 -->
          <div class="bg-gray-50 rounded-xl p-5">
            <div class="flex justify-between items-center mb-3">
              <span class="text-sm font-medium text-gray-700">光照阈值</span>
              <span class="text-sm font-bold text-blue-600">{{ threshold }} lux</span>
            </div>
            <input
              v-model.number="threshold"
              type="range"
              min="10"
              max="500"
              step="10"
              :disabled="device.status === 'offline'"
              :class="['w-full h-2 bg-gray-200 rounded-lg appearance-none', device.status === 'offline' ? 'cursor-not-allowed opacity-50' : 'cursor-pointer']"
            />
            <button
              @click="handleSetThreshold"
              :disabled="loading || device.status === 'offline'"
              :class="['mt-3 w-full bg-orange-500 text-white py-2 rounded-xl hover:bg-orange-600 transition-all font-semibold',
                (loading || device.status === 'offline') ? 'opacity-50 cursor-not-allowed' : '']"
            >
              设置阈值
            </button>
          </div>

          <!-- 命令执行结果 -->
          <div v-if="commandResult" :class="['rounded-xl p-4', commandResult.success ? 'bg-green-50' : 'bg-red-50']">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium" :class="commandResult.success ? 'text-green-700' : 'text-red-700'">
                {{ commandResult.success ? '✓ 执行成功' : '✗ 执行失败' }}
              </span>
              <span class="text-xs text-gray-500">{{ commandResult.ts }}</span>
            </div>
            <p class="text-sm text-gray-600">{{ commandResult.message }}</p>
          </div>

          <!-- 设备日志 -->
          <div v-if="logs.length > 0" class="bg-gray-50 rounded-xl p-4">
            <div class="flex justify-between items-center mb-3">
              <span class="text-sm font-medium text-gray-700">设备日志</span>
              <button @click="logs = []" class="text-xs text-gray-500 hover:text-gray-700">清空</button>
            </div>
            <div class="max-h-48 overflow-y-auto space-y-2">
              <div v-for="(log, index) in logs" :key="index" class="flex items-start gap-2 text-xs">
                <span class="text-gray-400 flex-shrink-0">{{ formatTimestamp(log.ts) }}</span>
                <span class="px-2 py-0.5 rounded bg-blue-100 text-blue-700 text-xs font-medium">{{ log.action }}</span>
                <span class="text-gray-600">{{ log.detail }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 通用设备界面 -->
        <div v-else class="space-y-6">
          <!-- 状态卡片 -->
          <div
            :class="[
              'rounded-xl p-5 text-center',
              device.status === 'online' ? 'bg-green-50' : 'bg-red-50'
            ]"
          >
            <div class="text-sm text-gray-600 mb-1">设备状态</div>
            <div
              :class="[
                'text-2xl font-bold',
                device.status === 'online' ? 'text-green-600' : 'text-red-600'
              ]"
            >
              {{ device.status === 'online' ? '在线' : '离线' }}
            </div>
          </div>

          <!-- 元数据展示 -->
          <div v-if="device.metadata && Object.keys(device.metadata || {}).length > 0">
            <h4 class="font-semibold text-gray-700 mb-3">设备信息</h4>
            <div class="grid grid-cols-2 gap-3">
              <div
                v-for="(value, key) in (device.metadata || {})"
                :key="key"
                class="bg-gray-50 rounded-lg p-4"
              >
                <div class="text-xs text-gray-500 uppercase font-medium">{{ key }}</div>
                <div class="text-sm font-semibold text-gray-800">{{ value }}</div>
              </div>
            </div>
          </div>

          <!-- 离线提示 -->
          <div v-if="device.status === 'offline'" class="bg-orange-50 border-2 border-orange-200 rounded-xl p-4">
            <div class="flex items-center gap-2">
              <span class="text-xl">⚠️</span>
              <span class="text-orange-700 font-medium">设备当前离线，无法执行操作</span>
            </div>
          </div>

          <!-- 自定义命令 -->
          <div>
            <h4 class="font-semibold text-gray-700 mb-3">发送命令</h4>
            <div class="flex gap-3">
              <input
                v-model="customCommand"
                type="text"
                :disabled="device.status === 'offline'"
                :class="['flex-1 px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent', device.status === 'offline' ? 'opacity-50 cursor-not-allowed' : '']"
                placeholder="输入命令"
                @keyup.enter="sendCustomCommand"
              />
              <button
                @click="sendCustomCommand"
                :disabled="loading || !customCommand || device.status === 'offline'"
                :class="['px-6 py-3 bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-all font-semibold', (loading || !customCommand || device.status === 'offline') ? 'opacity-50 cursor-not-allowed' : '']"
              >
                发送
              </button>
            </div>
          </div>
        </div>

        <!-- 删除按钮（不受离线透明度影响，始终正常显示） -->
        <div class="pt-4 border-t border-gray-100 [&_*]:opacity-100">
          <button
            @click="handleDelete"
            :disabled="loading"
            class="w-full bg-red-500 text-white py-3 rounded-xl hover:bg-red-600 transition-all font-semibold disabled:opacity-50 disabled:cursor-not-allowed"
          >
            删除设备
          </button>
        </div>
      </div>
    </div>

    <!-- OTA升级模态框 -->
    <div
      v-if="showOtaModal"
      class="fixed inset-0 bg-black/50 z-60 flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-xl font-bold text-gray-800">⬆️ 远程 OTA 升级</h3>
          <button
            @click="showOtaModal = false"
            class="text-gray-400 hover:text-gray-600 text-2xl"
          >
            ×
          </button>
        </div>

        <div class="space-y-4">
          <div>
            <div class="flex gap-3 mb-4">
              <button
                @click="otaWithConfigUrl"
                :disabled="loading"
                class="flex-1 py-3 bg-cyan-500 text-white rounded-xl hover:bg-cyan-600 transition-all font-semibold disabled:opacity-50"
              >
                使用配置 URL
              </button>
              <button
                @click="otaWithCustomUrl"
                :disabled="loading || !customOtaUrl"
                class="flex-1 py-3 bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-all font-semibold disabled:opacity-50"
              >
                使用自定义 URL
              </button>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                自定义固件 URL (可选)
              </label>
              <input
                v-model="customOtaUrl"
                type="url"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="http://example.com/firmware.bin"
              />
            </div>
          </div>

          <div class="text-sm text-gray-500 bg-gray-50 p-4 rounded-xl">
            <p class="font-medium mb-2">⚠️ 升级注意事项：</p>
            <ul class="list-disc pl-5 space-y-1">
              <li>升级过程中请勿断电</li>
              <li>设备会自动重启完成升级</li>
              <li>可在设备日志中查看升级进度</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import store from '../store'
import api from '../services/api'

const props = defineProps({
  device: Object,
})

const emit = defineEmits(['close'])

const threshold = ref(100)
const customCommand = ref('')
const loading = ref(false)
const logs = ref([])
const commandResult = ref(null)
const showOtaModal = ref(false)
const customOtaUrl = ref('')

const icon = computed(() => store.getDeviceIcon(props.device?.type))

const lightOn = computed(() => {
  const meta = props.device?.metadata
  if (!meta) return false
  return meta.light === 'ON' || meta.light === true
})

const isAutoMode = computed(() => {
  const meta = props.device?.metadata
  if (!meta) return false
  return meta.mode === 'auto' || meta.mode === true
})

watch(
  () => props.device,
  (newDevice) => {
    if (newDevice?.metadata?.threshold) {
      threshold.value = newDevice.metadata.threshold
    }
  },
  { immediate: true }
)

const formatUptime = (seconds) => store.formatUptime(seconds)

const close = () => {
  emit('close')
}

const sendCommand = async (command) => {
  if (props.device.status === 'offline') {
    alert('设备离线，无法执行命令')
    return
  }
  
  const meta = props.device?.metadata
  const cmd = command.toUpperCase()

  let prevLight, prevMode
  if (meta) {
    prevLight = meta.light
    prevMode = meta.mode

    if (cmd === 'ON') {
      meta.light = 'ON'
      meta.mode = 'manual'
    } else if (cmd === 'OFF') {
      meta.light = 'OFF'
      meta.mode = 'manual'
    } else if (cmd === 'AUTO') {
      meta.mode = 'auto'
    }
  }

  loading.value = true
  try {
    await api.sendCommand(props.device.id, command)
  } catch (e) {
    if (meta) {
      meta.light = prevLight
      meta.mode = prevMode
    }
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const sendCustomCommand = async () => {
  if (!customCommand.value) return
  await sendCommand(customCommand.value)
  customCommand.value = ''
}

const fetchLogs = async () => {
  if (props.device.status === 'offline') {
    alert('设备离线，无法获取日志')
    return
  }
  
  loading.value = true
  try {
    const response = await api.sendCommand(props.device.id, 'LOGS')
    if (response && response.logs) {
      logs.value = response.logs
    }
  } catch (e) {
    alert('获取日志失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

const formatTimestamp = (seconds) => {
  const date = new Date(seconds * 1000)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

const otaWithConfigUrl = async () => {
  showOtaModal.value = false
  await sendCommand('UPDATE')
}

const otaWithCustomUrl = async () => {
  if (!customOtaUrl.value) {
    alert('请先输入固件 URL')
    return
  }
  
  showOtaModal.value = false
  await sendCommand(`UPDATE=${customOtaUrl.value}`)
  customOtaUrl.value = ''
}

const handleSetThreshold = async () => {
  if (props.device.status === 'offline') {
    alert('设备离线，无法设置阈值')
    return
  }
  
  const meta = props.device?.metadata
  const prevThreshold = meta?.threshold

  if (meta) {
    meta.threshold = threshold.value
  }

  loading.value = true
  try {
    await api.setThreshold(threshold.value)
  } catch (e) {
    if (meta) {
      meta.threshold = prevThreshold
    }
    alert(e.message)
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('确定删除该设备吗？')) return
  
  loading.value = true
  try {
    await api.deleteDevice(props.device.id)
    await store.loadDevices()
    await store.loadStats()
    close()
  } catch (e) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}

// 事件处理函数
const handleCommandResult = (event) => {
  const result = event.detail
  if (result) {
    commandResult.value = result
    // 3秒后自动清除结果
    setTimeout(() => {
      if (commandResult.value === result) {
        commandResult.value = null
      }
    }, 5000)
  }
}

const handleDeviceLogs = (event) => {
  const data = event.detail
  if (data && data.logs) {
    logs.value = data.logs
  }
}

const handleDeviceError = (event) => {
  const error = event.detail
  if (error) {
    alert(`设备错误: ${error.detail || '未知错误'}`)
  }
}

// 生命周期钩子
onMounted(() => {
  window.addEventListener('commandResult', handleCommandResult)
  window.addEventListener('deviceLogs', handleDeviceLogs)
  window.addEventListener('deviceError', handleDeviceError)
})

onUnmounted(() => {
  window.removeEventListener('commandResult', handleCommandResult)
  window.removeEventListener('deviceLogs', handleDeviceLogs)
  window.removeEventListener('deviceError', handleDeviceError)
})
</script>
