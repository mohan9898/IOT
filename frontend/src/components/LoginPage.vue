<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
    <div class="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md transform transition-all">
      <div class="text-center mb-8">
        <div class="text-5xl mb-4">🔌</div>
        <h1 class="text-3xl font-bold text-gray-800">元枢智能物联系统</h1>
        <p class="text-gray-500 mt-2">太一语音助手 · 智能物联</p>
      </div>

      <!-- 登录表单 -->
      <div v-if="!showRegister" class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">用户名</label>
          <input
            v-model="loginUsername"
            type="text"
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
            placeholder="请输入用户名"
            @keyup.enter="handleLogin"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">密码</label>
          <input
            v-model="loginPassword"
            type="password"
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
            placeholder="请输入密码"
            @keyup.enter="handleLogin"
          />
        </div>
        <button
          @click="handleLogin"
          :disabled="loading"
          class="w-full bg-blue-500 text-white py-3 rounded-xl hover:bg-blue-600 transition-all font-semibold disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="loading" class="flex items-center justify-center">
            <span class="spinner mr-2" style="width: 20px; height: 20px; border-width: 2px;"></span>
            登录中...
          </span>
          <span v-else>登录</span>
        </button>
        <p class="text-center text-gray-500">
          还没有账号？
          <button @click="showRegister = true" class="text-blue-500 hover:underline font-medium">立即注册</button>
        </p>
      </div>

      <!-- 注册表单 -->
      <div v-else class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">用户名</label>
          <input
            v-model="registerUsername"
            type="text"
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
            placeholder="至少3个字符"
            @keyup.enter="handleRegister"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">密码</label>
          <input
            v-model="registerPassword"
            type="password"
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
            placeholder="至少6个字符"
            @keyup.enter="handleRegister"
          />
        </div>
        <button
          @click="handleRegister"
          :disabled="loading"
          class="w-full bg-green-500 text-white py-3 rounded-xl hover:bg-green-600 transition-all font-semibold disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="loading" class="flex items-center justify-center">
            <span class="spinner mr-2" style="width: 20px; height: 20px; border-width: 2px;"></span>
            注册中...
          </span>
          <span v-else>注册</span>
        </button>
        <p class="text-center text-gray-500">
          已有账号？
          <button @click="showRegister = false" class="text-blue-500 hover:underline font-medium">立即登录</button>
        </p>
      </div>

      <!-- 错误提示 -->
      <transition name="fade">
        <p v-if="error" class="mt-6 text-red-500 text-center text-sm bg-red-50 py-3 rounded-xl">
          {{ error }}
        </p>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import store from '../store'

const showRegister = ref(false)
const loginUsername = ref('')
const loginPassword = ref('')
const registerUsername = ref('')
const registerPassword = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!loginUsername.value || !loginPassword.value) {
    error.value = '请输入用户名和密码'
    return
  }
  
  loading.value = true
  error.value = ''
  
  try {
    await store.login(loginUsername.value, loginPassword.value)
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  if (registerUsername.value.length < 3) {
    error.value = '用户名至少需要3个字符'
    return
  }
  if (registerPassword.value.length < 6) {
    error.value = '密码至少需要6个字符'
    return
  }
  
  loading.value = true
  error.value = ''
  
  try {
    await store.register(registerUsername.value, registerPassword.value)
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>
