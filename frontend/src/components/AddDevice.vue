<template>
  <div class="bg-white rounded-2xl shadow-lg p-8">
    <h2 class="text-xl font-bold text-gray-800 mb-6 flex items-center">
      <span class="text-2xl mr-2">➕</span>
      添加新设备
    </h2>
    
    <form @submit.prevent="handleAddDevice" class="space-y-6">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">设备 ID</label>
          <input
            v-model="form.id"
            type="text"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
            placeholder="设备唯一标识"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">设备名称</label>
          <input
            v-model="form.name"
            type="text"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
            placeholder="设备名称"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">设备类型</label>
          <select
            v-model="form.type"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          >
            <option value="">请选择设备类型</option>
            <option v-for="dt in deviceTypes" :key="dt.type_id" :value="dt.type_id">
              {{ dt.icon }} {{ dt.name }}
            </option>
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">备注信息</label>
          <input
            v-model="form.note"
            type="text"
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
            placeholder="可选备注"
          />
        </div>
      </div>
      
      <button
        type="submit"
        :disabled="loading"
        class="w-full md:w-auto px-8 py-3 bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-all font-semibold disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
      >
        <span v-if="loading" class="flex items-center">
          <span class="spinner mr-2" style="width: 20px; height: 20px; border-width: 2px;"></span>
          添加中...
        </span>
        <span v-else>添加设备</span>
      </button>
      
      <transition name="fade">
        <p v-if="error" class="text-red-500 text-sm bg-red-50 py-3 px-4 rounded-xl">
          {{ error }}
        </p>
        <p v-if="success" class="text-green-500 text-sm bg-green-50 py-3 px-4 rounded-xl">
          {{ success }}
        </p>
      </transition>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import store from '../store'
import api from '../services/api'

const form = ref({
  id: '',
  name: '',
  type: '',
  note: '',
})

const loading = ref(false)
const error = ref('')
const success = ref('')

const deviceTypes = computed(() => store.deviceTypes)

const handleAddDevice = async () => {
  if (!form.value.id || !form.value.name || !form.value.type) {
    error.value = '请填写完整信息'
    return
  }
  
  loading.value = true
  error.value = ''
  success.value = ''
  
  try {
    await api.addDevice(form.value.id, form.value.name, form.value.type, form.value.note)
    success.value = '设备添加成功！'
    
    // 重置表单
    form.value = { id: '', name: '', type: '', note: '' }
    
    // 刷新数据
    await store.loadDevices()
    await store.loadStats()
    
    // 3秒后清除成功消息
    setTimeout(() => {
      success.value = ''
    }, 3000)
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await store.loadDeviceTypes()
})
</script>
