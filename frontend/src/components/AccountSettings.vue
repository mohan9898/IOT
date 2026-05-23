<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
    <!-- 修改用户名 -->
    <div class="bg-white rounded-2xl shadow-lg p-8">
      <h2 class="text-xl font-bold text-gray-800 mb-6 flex items-center">
        <span class="text-2xl mr-2">👤</span>
        修改用户名
      </h2>
      
      <form @submit.prevent="handleUpdateAccount" class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">当前用户名</label>
          <input
            v-model="accountForm.oldUsername"
            type="text"
            disabled
            class="w-full px-4 py-3 border border-gray-200 rounded-xl bg-gray-50 text-gray-500"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">新用户名</label>
          <input
            v-model="accountForm.newUsername"
            type="text"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
            placeholder="至少3个字符"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">确认密码</label>
          <input
            v-model="accountForm.password"
            type="password"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
            placeholder="输入密码以确认"
          />
        </div>
        
        <button
          type="submit"
          :disabled="accountLoading"
          class="w-full bg-blue-500 text-white py-3 rounded-xl hover:bg-blue-600 transition-all font-semibold disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
        >
          <span v-if="accountLoading" class="flex items-center">
            <span class="spinner mr-2" style="width: 20px; height: 20px; border-width: 2px;"></span>
            修改中...
          </span>
          <span v-else>修改用户名</span>
        </button>
        
        <transition name="fade">
          <p v-if="accountError" class="text-red-500 text-sm bg-red-50 py-3 px-4 rounded-xl">
            {{ accountError }}
          </p>
          <p v-if="accountSuccess" class="text-green-500 text-sm bg-green-50 py-3 px-4 rounded-xl">
            {{ accountSuccess }}
          </p>
        </transition>
      </form>
    </div>

    <!-- 修改密码 -->
    <div class="bg-white rounded-2xl shadow-lg p-8">
      <h2 class="text-xl font-bold text-gray-800 mb-6 flex items-center">
        <span class="text-2xl mr-2">🔐</span>
        修改密码
      </h2>
      
      <form @submit.prevent="handleUpdatePassword" class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">用户名</label>
          <input
            v-model="passwordForm.username"
            type="text"
            disabled
            class="w-full px-4 py-3 border border-gray-200 rounded-xl bg-gray-50 text-gray-500"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">当前密码</label>
          <input
            v-model="passwordForm.oldPassword"
            type="password"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-all"
            placeholder="输入当前密码"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">新密码</label>
          <input
            v-model="passwordForm.newPassword"
            type="password"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-all"
            placeholder="至少6个字符"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">确认新密码</label>
          <input
            v-model="passwordForm.confirmPassword"
            type="password"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-all"
            placeholder="再次输入新密码"
          />
        </div>
        
        <button
          type="submit"
          :disabled="passwordLoading"
          class="w-full bg-orange-500 text-white py-3 rounded-xl hover:bg-orange-600 transition-all font-semibold disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
        >
          <span v-if="passwordLoading" class="flex items-center">
            <span class="spinner mr-2" style="width: 20px; height: 20px; border-width: 2px;"></span>
            修改中...
          </span>
          <span v-else>修改密码</span>
        </button>
        
        <transition name="fade">
          <p v-if="passwordError" class="text-red-500 text-sm bg-red-50 py-3 px-4 rounded-xl">
            {{ passwordError }}
          </p>
          <p v-if="passwordSuccess" class="text-green-500 text-sm bg-green-50 py-3 px-4 rounded-xl">
            {{ passwordSuccess }}
          </p>
        </transition>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import store from '../store'
import api from '../services/api'

const accountForm = ref({
  oldUsername: '',
  newUsername: '',
  password: '',
})

const passwordForm = ref({
  username: '',
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const accountLoading = ref(false)
const passwordLoading = ref(false)
const accountError = ref('')
const accountSuccess = ref('')
const passwordError = ref('')
const passwordSuccess = ref('')

const handleUpdateAccount = async () => {
  if (accountForm.value.newUsername.length < 3) {
    accountError.value = '用户名至少需要3个字符'
    return
  }
  if (!accountForm.value.password) {
    accountError.value = '请输入密码以确认'
    return
  }
  
  accountLoading.value = true
  accountError.value = ''
  accountSuccess.value = ''
  
  try {
    const result = await api.updateAccount(
      accountForm.value.oldUsername,
      accountForm.value.newUsername,
      accountForm.value.password
    )
    
    accountSuccess.value = '用户名修改成功！'
    store.currentUsername = result.username
    accountForm.value.oldUsername = result.username
    passwordForm.value.username = result.username
    accountForm.value.password = ''
    
    setTimeout(() => {
      accountSuccess.value = ''
    }, 3000)
  } catch (e) {
    accountError.value = e.message
  } finally {
    accountLoading.value = false
  }
}

const handleUpdatePassword = async () => {
  if (!passwordForm.value.oldPassword) {
    passwordError.value = '请输入当前密码'
    return
  }
  if (passwordForm.value.newPassword.length < 6) {
    passwordError.value = '新密码至少需要6个字符'
    return
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    passwordError.value = '两次输入的新密码不一致'
    return
  }
  
  passwordLoading.value = true
  passwordError.value = ''
  passwordSuccess.value = ''
  
  try {
    await api.updatePassword(
      passwordForm.value.username,
      passwordForm.value.oldPassword,
      passwordForm.value.newPassword
    )
    
    passwordSuccess.value = '密码修改成功！'
    passwordForm.value.oldPassword = ''
    passwordForm.value.newPassword = ''
    passwordForm.value.confirmPassword = ''
    
    setTimeout(() => {
      passwordSuccess.value = ''
    }, 3000)
  } catch (e) {
    passwordError.value = e.message
  } finally {
    passwordLoading.value = false
  }
}

onMounted(() => {
  accountForm.value.oldUsername = store.currentUsername
  passwordForm.value.username = store.currentUsername
})
</script>
