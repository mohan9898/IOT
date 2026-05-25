<template>
  <div class="space-y-5">
    <!-- 统计卡片 -->
    <div class="grid grid-cols-3 gap-4">
      <div class="bg-white rounded-xl p-4 shadow-sm border border-gray-100">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ stats.total }}</div>
            <div class="text-xs text-gray-500 mt-0.5 whitespace-nowrap">总设备数</div>
          </div>
          <div class="text-2xl flex-shrink-0">📊</div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl p-4 shadow-sm border border-gray-100">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-green-500">{{ stats.online }}</div>
            <div class="text-xs text-gray-500 mt-0.5 whitespace-nowrap">在线设备</div>
          </div>
          <div class="text-2xl flex-shrink-0">🟢</div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl p-4 shadow-sm border border-gray-100">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-2xl font-bold text-red-500">{{ stats.offline }}</div>
            <div class="text-xs text-gray-500 mt-0.5 whitespace-nowrap">离线设备</div>
          </div>
          <div class="text-2xl flex-shrink-0">🔴</div>
        </div>
      </div>
    </div>

    <!-- 设备列表 -->
    <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-bold text-gray-800 flex items-center whitespace-nowrap">
          <span class="text-xl mr-2">📱</span>
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
      
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
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
