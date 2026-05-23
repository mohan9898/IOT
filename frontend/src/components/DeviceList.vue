<template>
  <div class="space-y-6">
    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-4xl font-bold text-gray-800">{{ stats.total }}</div>
            <div class="text-gray-500 mt-1 font-medium">总设备数</div>
          </div>
          <div class="text-4xl">📊</div>
        </div>
      </div>
      
      <div class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-4xl font-bold text-green-500">{{ stats.online }}</div>
            <div class="text-gray-500 mt-1 font-medium">在线设备</div>
          </div>
          <div class="text-4xl">🟢</div>
        </div>
      </div>
      
      <div class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-4xl font-bold text-red-500">{{ stats.offline }}</div>
            <div class="text-gray-500 mt-1 font-medium">离线设备</div>
          </div>
          <div class="text-4xl">🔴</div>
        </div>
      </div>
    </div>

    <!-- 设备列表 -->
    <div class="bg-white rounded-2xl shadow-lg p-8">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold text-gray-800 flex items-center">
          <span class="text-2xl mr-2">📱</span>
          设备列表
        </h2>
        <button
          @click="refresh"
          :disabled="store.loading.devices"
          class="px-4 py-2 text-blue-500 hover:bg-blue-50 rounded-xl transition-colors flex items-center disabled:opacity-50"
        >
          <span v-if="store.loading.devices" class="spinner mr-2" style="width: 16px; height: 16px; border-width: 2px;"></span>
          <span v-else class="mr-2">🔄</span>
          刷新
        </button>
      </div>
      
      <div v-if="store.loading.devices && devices.length === 0" class="flex justify-center py-16">
        <div class="spinner"></div>
      </div>
      
      <div v-else-if="devices.length === 0" class="text-center py-16 text-gray-500">
        <div class="text-6xl mb-4">📭</div>
        <p class="text-lg">暂无设备，点击"添加设备"开始添加</p>
      </div>
      
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <DeviceCard
          v-for="device in devices"
          :key="device.id"
          :device="device"
          @select="selectDevice"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import store from '../store'
import DeviceCard from './DeviceCard.vue'

const devices = computed(() => store.devices)
const stats = computed(() => store.stats)

const selectDevice = (device) => {
  store.selectDevice(device)
}

const refresh = async () => {
  await Promise.all([store.loadDevices(), store.loadStats()])
}

onMounted(async () => {
  await Promise.all([store.loadDevices(), store.loadStats()])
})
</script>
