<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-xl sm:text-2xl font-bold text-gray-800 whitespace-nowrap">控制记录</h1>
      <button
        @click="loadRecords"
        :disabled="loading"
        class="px-4 py-2 text-sm bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-colors disabled:opacity-50"
      >
        {{ loading ? '加载中...' : '刷新' }}
      </button>
    </div>

    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div class="bg-white rounded-xl p-4 shadow-sm border border-gray-100">
        <p class="text-xs text-gray-500">总记录</p>
        <p class="text-2xl font-bold text-gray-800">{{ stats.total }}</p>
      </div>
      <div class="bg-white rounded-xl p-4 shadow-sm border border-gray-100">
        <p class="text-xs text-gray-500">今日</p>
        <p class="text-2xl font-bold text-blue-600">{{ stats.today }}</p>
      </div>
      <div class="bg-white rounded-xl p-4 shadow-sm border border-gray-100">
        <p class="text-xs text-gray-500">成功</p>
        <p class="text-2xl font-bold text-green-600">{{ stats.success }}</p>
      </div>
      <div class="bg-white rounded-xl p-4 shadow-sm border border-gray-100">
        <p class="text-xs text-gray-500">失败</p>
        <p class="text-2xl font-bold text-red-500">{{ stats.failed }}</p>
      </div>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      <div class="p-4 border-b border-gray-100 flex flex-wrap gap-3 items-center">
        <div class="flex items-center gap-2">
          <span class="text-sm text-gray-500">设备筛选：</span>
          <select
            v-model="filterDeviceId"
            @change="loadRecords"
            class="px-3 py-2 border border-gray-300 rounded-xl text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="">全部设备</option>
            <option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }}</option>
          </select>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table v-if="records.length > 0" class="w-full">
          <thead>
            <tr class="bg-gray-50 text-left">
              <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase whitespace-nowrap">设备</th>
              <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase whitespace-nowrap">命令</th>
              <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase whitespace-nowrap">状态</th>
              <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase hidden sm:table-cell whitespace-nowrap">参数</th>
              <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase hidden md:table-cell whitespace-nowrap">发送时间</th>
              <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase hidden md:table-cell whitespace-nowrap">响应</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="rec in records" :key="rec.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-4 py-3.5 max-w-[160px]">
                <div class="font-medium text-gray-800 text-sm truncate">{{ getDeviceName(rec.device_id) }}</div>
                <div class="text-xs text-gray-400 font-mono truncate">{{ rec.device_id }}</div>
              </td>
              <td class="px-4 py-3.5">
                <span class="inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-semibold whitespace-nowrap"
                  :class="commandClass(rec.command)"
                >{{ rec.command }}</span>
              </td>
              <td class="px-4 py-3.5">
                <span :class="['inline-flex items-center gap-1 text-xs font-semibold whitespace-nowrap', statusClass(rec.status)]">
                  <span :class="['w-2 h-2 rounded-full', statusDotClass(rec.status)]"></span>
                  {{ rec.status || 'pending' }}
                </span>
              </td>
              <td class="px-4 py-3.5 hidden sm:table-cell">
                <span class="text-xs text-gray-500 whitespace-nowrap">{{ formatParams(rec.parameters) }}</span>
              </td>
              <td class="px-4 py-3.5 hidden md:table-cell">
                <span class="text-xs text-gray-500 whitespace-nowrap">{{ formatTime(rec.sent_at) }}</span>
              </td>
              <td class="px-4 py-3.5 hidden md:table-cell">
                <span class="text-xs text-gray-400 max-w-[150px] truncate block">{{ rec.response || '-' }}</span>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-else class="p-12 text-center text-gray-400">
          <p class="text-4xl mb-3">📋</p>
          <p class="text-sm">暂无控制记录</p>
          <p class="text-xs mt-1">发送控制命令后将在此处显示</p>
        </div>
      </div>

      <div v-if="totalPages > 1" class="flex justify-between items-center p-4 border-t border-gray-100">
        <span class="text-sm text-gray-500">共 {{ total }} 条记录</span>
        <div class="flex gap-1">
          <button
            @click="goPage(page - 1)"
            :disabled="page <= 1"
            class="px-3 py-1.5 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
          >上一页</button>
          <button
            v-for="p in visiblePages"
            :key="p"
            @click="goPage(p)"
            :class="['px-3 py-1.5 text-sm rounded-lg', p === page ? 'bg-blue-500 text-white' : 'border border-gray-300 hover:bg-gray-50']"
          >{{ p }}</button>
          <button
            @click="goPage(page + 1)"
            :disabled="page >= totalPages"
            class="px-3 py-1.5 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
          >下一页</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import store from '../store'

const records = ref([])
const stats = ref({ total: 0, today: 0, success: 0, failed: 0 })
const loading = ref(false)
const page = ref(1)
const total = ref(0)
const totalPages = ref(0)
const filterDeviceId = ref('')
const devices = ref([])

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, page.value - 2)
  const end = Math.min(totalPages.value, page.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

async function loadRecords() {
  loading.value = true
  try {
    const data = await store.api.getControlRecords(filterDeviceId.value, page.value, 20)
    records.value = data.records || []
    total.value = data.total || 0
    totalPages.value = data.total_pages || 0
  } catch (e) {
    console.error('加载控制记录失败:', e)
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  try {
    const data = await store.api.getControlStats()
    stats.value = data
  } catch (e) {
    console.error('加载统计失败:', e)
  }
}

async function loadDevices() {
  try {
    const data = await store.api.getDevices()
    devices.value = Array.isArray(data) ? data : []
  } catch (e) {
    devices.value = []
  }
}

function goPage(p) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  loadRecords()
}

function getDeviceName(id) {
  const d = devices.value.find(d => d.id === id)
  return d ? d.name : id
}

function commandClass(cmd) {
  const c = (cmd || '').toUpperCase()
  if (c === 'ON') return 'bg-green-100 text-green-700'
  if (c === 'OFF') return 'bg-gray-200 text-gray-700'
  if (c === 'AUTO') return 'bg-blue-100 text-blue-700'
  return 'bg-purple-100 text-purple-700'
}

function statusClass(status) {
  if (status === 'success') return 'text-green-600'
  if (status === 'failed') return 'text-red-500'
  return 'text-gray-400'
}

function statusDotClass(status) {
  if (status === 'success') return 'bg-green-500'
  if (status === 'failed') return 'bg-red-500'
  return 'bg-gray-300'
}

function formatParams(params) {
  if (!params || Object.keys(params).length === 0) return '-'
  try {
    return Object.entries(params).map(([k, v]) => `${k}:${v}`).join(', ')
  } catch {
    return '-'
  }
}

function formatTime(ts) {
  if (!ts) return '-'
  try {
    const d = new Date(ts)
    return d.toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' })
  } catch {
    return ts
  }
}

onMounted(() => {
  loadDevices()
  loadStats()
  loadRecords()
})
</script>