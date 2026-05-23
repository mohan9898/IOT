<template>
  <div v-if="device" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" @click.self="close">
    <div class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
      <!-- 头部 -->
      <div class="p-6 border-b border-gray-100">
        <div class="flex justify-between items-start">
          <div class="flex items-center">
            <span class="text-5xl mr-4">{{ icon }}</span>
            <div>
              <h2 class="text-2xl font-bold text-gray-800">{{ device.name }}</h2>
              <p class="text-sm text-gray-500 font-mono">{{ device.id }}</p>
            </div>
          </div>
          <button @click="close" class="text-gray-400 hover:text-gray-600 text-2xl p-2 hover:bg-gray-100 rounded-full transition-colors">
            ✕
          </button>
        </div>
      </div>

      <!-- 内容 -->
      <div class="p-6 space-y-6">
        <!-- 智能灯特殊界面 -->
        <div v-if="device.type === 'smart_light'" class="space-y-6">
          <!-- 传感器数据 -->
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-gradient-to-br from-blue-50 to-blue-100 rounded-xl p-5">
              <div class="text-sm text-blue-600 mb-1 font-medium">光照强度</div>
              <div class="text-3xl font-bold text-blue-700">{{ device.metadata?.lux || 0 }} lux</div>
            </div>
            <div class="bg-gradient-to-br from-green-50 to-green-100 rounded-xl p-5">
              <div class="text-sm text-green-600 mb-1 font-medium">信号强度</div>
              <div class="text-3xl font-bold text-green-700">{{ device.metadata?.rssi || 0 }} dBm</div>
            </div>
            <div class="bg-gradient-to-br from-purple-50 to-purple-100 rounded-xl p-5">
              <div class="text-sm text-purple-600 mb-1 font-medium">人体感应</div>
              <div :class="['text-3xl font-bold', device.metadata?.presence ? 'text-purple-700' : 'text-gray-400']">
                {{ device.metadata?.presence ? '检测到' : '无' }}
              </div>
            </div>
            <div class="bg-gradient-to-br from-orange-50 to-orange-100 rounded-xl p-5">
              <div class="text-sm text-orange-600 mb-1 font-medium">运行时长</div>
              <div class="text-3xl font-bold text-orange-700">{{ formatUptime(device.metadata?.uptime || 0) }}</div>
            </div>
          </div>

          <!-- 快捷控制 -->
          <div class="grid grid-cols-3 gap-3">
            <button
              @click="sendCommand('ON')"
              :disabled="loading"
              class="py-4 bg-green-500 text-white rounded-xl hover:bg-green-600 transition-all font-semibold disabled:opacity-50"
            >
              🔆 开灯
            </button>
            <button
              @click="sendCommand('OFF')"
              :disabled="loading"
              class="py-4 bg-gray-500 text-white rounded-xl hover:bg-gray-600 transition-all font-semibold disabled:opacity-50"
            >
              🌙 关灯
            </button>
            <button
              @click="sendCommand('AUTO')"
              :disabled="loading"
              class="py-4 bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-all font-semibold disabled:opacity-50"
            >
              🤖 自动
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
              class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
            />
            <button
              @click="handleSetThreshold"
              :disabled="loading"
              class="mt-3 w-full bg-orange-500 text-white py-2 rounded-xl hover:bg-orange-600 transition-all font-semibold disabled:opacity-50"
            >
              设置阈值
            </button>
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

          <!-- 自定义命令 -->
          <div>
            <h4 class="font-semibold text-gray-700 mb-3">发送命令</h4>
            <div class="flex gap-3">
              <input
                v-model="customCommand"
                type="text"
                class="flex-1 px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="输入命令"
                @keyup.enter="sendCustomCommand"
              />
              <button
                @click="sendCustomCommand"
                :disabled="loading || !customCommand"
                class="px-6 py-3 bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-all font-semibold disabled:opacity-50"
              >
                发送
              </button>
            </div>
          </div>
        </div>

        <!-- 删除按钮 -->
        <div class="pt-4 border-t border-gray-100">
          <button
            @click="handleDelete"
            :disabled="loading"
            class="w-full bg-red-500 text-white py-3 rounded-xl hover:bg-red-600 transition-all font-semibold disabled:opacity-50"
          >
            删除设备
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import store from '../store'
import api from '../services/api'

const props = defineProps({
  device: Object,
})

const emit = defineEmits(['close'])

const threshold = ref(100)
const customCommand = ref('')
const loading = ref(false)

const icon = computed(() => store.getDeviceIcon(props.device?.type))

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
  loading.value = true
  try {
    await api.sendCommand(props.device.id, command)
    await store.loadDevices()
  } catch (e) {
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

const handleSetThreshold = async () => {
  loading.value = true
  try {
    await api.setThreshold(threshold.value)
    await store.loadDevices()
  } catch (e) {
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
</script>
