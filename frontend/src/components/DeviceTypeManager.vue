<template>
  <div class="space-y-6">
    <!-- 添加设备类型 -->
    <div class="bg-white rounded-2xl shadow-lg p-8">
      <h2 class="text-xl font-bold text-gray-800 mb-6 flex items-center">
        <span class="text-2xl mr-2">📦</span>
        添加设备类型
      </h2>
      
      <form @submit.prevent="handleAddType" class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">类型标识</label>
            <input
              v-model="form.typeId"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
              placeholder="如: smart_switch"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">类型名称</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
              placeholder="如: 智能开关"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">图标</label>
            <input
              v-model="form.icon"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
              placeholder="📦"
            />
          </div>
          
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-2">描述</label>
            <textarea
              v-model="form.description"
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
              rows="2"
              placeholder="设备类型描述"
            ></textarea>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">支持的命令</label>
            <input
              v-model="form.commands"
              type="text"
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
              placeholder="ON, OFF, TOGGLE"
            />
          </div>
        </div>
        
        <button
          type="submit"
          :disabled="loading"
          class="px-8 py-3 bg-green-500 text-white rounded-xl hover:bg-green-600 transition-all font-semibold disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
        >
          <span v-if="loading" class="flex items-center">
            <span class="spinner mr-2" style="width: 20px; height: 20px; border-width: 2px;"></span>
            添加中...
          </span>
          <span v-else>添加类型</span>
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

    <!-- 设备类型列表 -->
    <div class="bg-white rounded-2xl shadow-lg p-8">
      <h2 class="text-xl font-bold text-gray-800 mb-6 flex items-center">
        <span class="text-2xl mr-2">📋</span>
        设备类型列表
      </h2>
      
      <div v-if="store.loading.deviceTypes" class="flex justify-center py-12">
        <div class="spinner"></div>
      </div>
      
      <div v-else-if="deviceTypes.length === 0" class="text-center py-12 text-gray-500">
        暂无设备类型
      </div>
      
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="dt in deviceTypes"
          :key="dt.type_id"
          class="border-2 border-gray-100 rounded-xl p-6 hover:shadow-md transition-all"
        >
          <div class="flex items-center mb-4">
            <span class="text-4xl mr-4">{{ dt.icon }}</span>
            <div class="flex-1 min-w-0">
              <h3 class="font-bold text-gray-800 truncate">{{ dt.name }}</h3>
              <p class="text-xs text-gray-500 font-mono truncate">{{ dt.type_id }}</p>
            </div>
          </div>
          
          <p class="text-sm text-gray-600 mb-4">{{ dt.description || '暂无描述' }}</p>
          
          <div v-if="dt.supported_commands && dt.supported_commands.length > 0">
            <div class="text-xs text-gray-500 mb-2 font-medium">支持命令:</div>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="cmd in parseCommands(dt.supported_commands)"
                :key="cmd"
                class="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-xs font-medium"
              >
                {{ cmd }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import store from '../store'
import api from '../services/api'

const form = ref({
  typeId: '',
  name: '',
  icon: '📦',
  description: '',
  commands: '',
})

const loading = ref(false)
const error = ref('')
const success = ref('')

const deviceTypes = computed(() => store.deviceTypes)

const parseCommands = (commands) => {
  if (!commands) return []
  if (Array.isArray(commands)) return commands
  try {
    return JSON.parse(commands)
  } catch {
    return []
  }
}

const handleAddType = async () => {
  if (!form.value.typeId || !form.value.name || !form.value.icon) {
    error.value = '请填写完整信息'
    return
  }
  
  loading.value = true
  error.value = ''
  success.value = ''
  
  try {
    const commands = form.value.commands
      ? form.value.commands.split(',').map((c) => c.trim())
      : []
    
    await api.addDeviceType(form.value.typeId, form.value.name, form.value.description, form.value.icon, commands)
    success.value = '设备类型添加成功！'
    
    // 重置表单
    form.value = { typeId: '', name: '', icon: '📦', description: '', commands: '' }
    
    // 刷新数据
    await store.loadDeviceTypes()
    
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
