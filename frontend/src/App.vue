<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 登录页面 -->
    <LoginPage v-if="!store.loggedIn" />
    
    <!-- 主应用 -->
    <div v-else>
      <!-- 欢迎引导弹窗 -->
      <div
        v-if="showWelcome"
        class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
        @click.self="closeWelcome"
      >
        <div class="bg-white rounded-2xl p-8 max-w-md w-full shadow-2xl">
          <div class="text-center mb-6">
            <p class="text-6xl mb-4">👋</p>
            <h2 class="text-2xl font-bold text-gray-800 mb-2">欢迎使用 IoT 设备管理器</h2>
            <p class="text-gray-500">快速上手，轻松管理您的智能设备</p>
          </div>
          <div class="space-y-3 mb-6">
            <div class="flex items-center gap-3 p-3 bg-blue-50 rounded-xl">
              <span class="text-2xl">📊</span>
              <div>
                <p class="font-medium text-gray-800">仪表盘总览</p>
                <p class="text-sm text-gray-500">一屏查看设备状态和统计</p>
              </div>
            </div>
            <div class="flex items-center gap-3 p-3 bg-green-50 rounded-xl">
              <span class="text-2xl">➕</span>
              <div>
                <p class="font-medium text-gray-800">添加设备</p>
                <p class="text-sm text-gray-500">手动添加或 MQTT 自动注册</p>
              </div>
            </div>
            <div class="flex items-center gap-3 p-3 bg-purple-50 rounded-xl">
              <span class="text-2xl">🔔</span>
              <div>
                <p class="font-medium text-gray-800">实时告警</p>
                <p class="text-sm text-gray-500">设备离线时浏览器推送通知</p>
              </div>
            </div>
          </div>
          <div class="flex gap-3">
            <button
              @click="closeWelcome"
              class="flex-1 px-5 py-3 bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-colors font-medium"
            >
              开始使用
            </button>
            <button
              @click="closeWelcomeAndNotify"
              class="flex-1 px-5 py-3 bg-green-500 text-white rounded-xl hover:bg-green-600 transition-colors font-medium"
            >
              开启通知并开始
            </button>
          </div>
        </div>
      </div>

      <!-- 顶部导航栏 -->
      <nav class="bg-white shadow-lg sticky top-0 z-40">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between items-center h-16">
            <!-- Logo -->
            <div class="flex items-center">
              <span class="text-3xl mr-3">🔌</span>
              <span class="text-lg sm:text-xl font-bold text-gray-800 hidden sm:inline">IoT 设备管理</span>
              <span class="text-lg font-bold text-gray-800 sm:hidden">IoT</span>
            </div>
            
            <!-- 桌面端导航 -->
            <div class="hidden md:flex items-center space-x-1">
              <button
                v-for="tab in tabs"
                :key="tab.id"
                @click="activeTab = tab.id"
                :class="[
                  'px-3 lg:px-5 py-2.5 rounded-xl font-medium transition-all duration-200 text-sm lg:text-base',
                  activeTab === tab.id
                    ? 'bg-blue-500 text-white shadow-md'
                    : 'text-gray-600 hover:bg-gray-100'
                ]"
              >
                <span class="mr-1 lg:mr-2">{{ tab.icon }}</span>
                {{ tab.label }}
              </button>
              
              <div class="w-px h-8 bg-gray-200 mx-2"></div>
              
              <!-- 通知开关 -->
              <button
                @click="toggleNotifications"
                :class="[
                  'px-3 py-2.5 rounded-xl font-medium transition-all text-sm',
                  store.notificationsEnabled
                    ? 'bg-green-100 text-green-700 hover:bg-green-200'
                    : 'bg-gray-100 text-gray-500 hover:bg-gray-200'
                ]"
                :title="store.notificationsEnabled ? '通知已开启' : '通知已关闭'"
              >
                <span class="text-lg">{{ store.notificationsEnabled ? '🔔' : '🔕' }}</span>
              </button>
              
              <div class="flex items-center mr-2">
                <span class="text-gray-500 mr-1">👤</span>
                <span class="text-gray-700 font-medium text-sm">{{ store.currentUsername }}</span>
              </div>
              
              <button
                @click="logout"
                class="px-3 py-2.5 text-red-500 hover:bg-red-50 rounded-xl font-medium text-sm transition-colors"
              >
                退出
              </button>
            </div>
            
            <!-- 移动端菜单按钮 -->
            <button
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="md:hidden p-2 rounded-xl hover:bg-gray-100 transition-colors"
            >
              <span class="text-2xl">{{ mobileMenuOpen ? '✕' : '☰' }}</span>
            </button>
          </div>
          
          <!-- 移动端菜单 -->
          <transition name="slide-fade">
            <div v-if="mobileMenuOpen" class="md:hidden pb-4">
              <div class="space-y-1">
                <button
                  v-for="tab in tabs"
                  :key="tab.id"
                  @click="activeTab = tab.id; mobileMenuOpen = false"
                  :class="[
                    'w-full text-left px-4 py-3 rounded-xl font-medium transition-all',
                    activeTab === tab.id
                      ? 'bg-blue-500 text-white'
                      : 'text-gray-600 hover:bg-gray-100'
                  ]"
                >
                  <span class="mr-3">{{ tab.icon }}</span>
                  {{ tab.label }}
                </button>
                
                <div class="pt-4 border-t border-gray-100">
                  <button
                    @click="toggleNotifications"
                    class="w-full text-left px-4 py-3 hover:bg-gray-100 rounded-xl font-medium"
                  >
                    <span class="mr-3">{{ store.notificationsEnabled ? '🔔' : '🔕' }}</span>
                    {{ store.notificationsEnabled ? '通知已开启' : '通知已关闭（点击开启）' }}
                  </button>
                  <div class="flex items-center px-4 py-3">
                    <span class="text-gray-500 mr-3">👤</span>
                    <span class="text-gray-700 font-medium">{{ store.currentUsername }}</span>
                  </div>
                  <button
                    @click="logout"
                    class="w-full text-left px-4 py-3 text-red-500 hover:bg-red-50 rounded-xl font-medium"
                  >
                    退出登录
                  </button>
                </div>
              </div>
            </div>
          </transition>
        </div>
      </nav>
      
      <!-- 主内容区 -->
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
        <transition name="fade" mode="out-in">
          <Dashboard v-if="activeTab === 'dashboard'" key="dashboard" @navigate="activeTab = $event" />
          <DeviceList v-else-if="activeTab === 'devices'" key="devices" />
          <AddDevice v-else-if="activeTab === 'add'" key="add" />
          <DeviceTypeManager v-else-if="activeTab === 'types'" key="types" />
          <AccountSettings v-else-if="activeTab === 'settings'" key="settings" />
        </transition>
      </main>
      
      <!-- 设备详情模态框 -->
      <DeviceDetailModal
        :device="store.selectedDevice"
        @close="store.closeDeviceModal()"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import store from './store'
import LoginPage from './components/LoginPage.vue'
import Dashboard from './components/Dashboard.vue'
import DeviceList from './components/DeviceList.vue'
import AddDevice from './components/AddDevice.vue'
import DeviceTypeManager from './components/DeviceTypeManager.vue'
import AccountSettings from './components/AccountSettings.vue'
import DeviceDetailModal from './components/DeviceDetailModal.vue'

const activeTab = ref('dashboard')
const mobileMenuOpen = ref(false)
const showWelcome = ref(false)

const tabs = computed(() => [
  { id: 'dashboard', label: '仪表盘', icon: '📊' },
  { id: 'devices', label: '设备列表', icon: '📱' },
  { id: 'add', label: '添加设备', icon: '➕' },
  { id: 'types', label: '设备类型', icon: '📦' },
  { id: 'settings', label: '账户设置', icon: '⚙️' },
])

function toggleNotifications() {
  if (store.notificationsEnabled) {
    store.disableNotifications()
  } else {
    store.enableNotifications()
  }
}

function closeWelcome() {
  showWelcome.value = false
  store.seenOnboarding = true
}

function closeWelcomeAndNotify() {
  store.enableNotifications()
  closeWelcome()
}

const logout = () => {
  if (confirm('确定要退出登录吗？')) {
    store.logout()
    activeTab.value = 'dashboard'
  }
}

// 检查是否显示欢迎引导
if (store.loggedIn && !store.seenOnboarding) {
  setTimeout(() => {
    showWelcome.value = true
  }, 500)
}
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}
.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(1, 0.5, 0.8, 1);
}
.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}
</style>