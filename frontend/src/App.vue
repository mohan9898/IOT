<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 登录页面 -->
    <LoginPage v-if="!store.loggedIn" />
    
    <!-- 主应用 -->
    <div v-else>
      <!-- 顶部导航栏 -->
      <nav class="bg-white shadow-lg sticky top-0 z-40">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between items-center h-16">
            <!-- Logo -->
            <div class="flex items-center">
              <span class="text-3xl mr-3">🔌</span>
              <span class="text-xl font-bold text-gray-800">IoT 设备管理</span>
            </div>
            
            <!-- 桌面端导航 -->
            <div class="hidden md:flex items-center space-x-2">
              <button
                v-for="tab in tabs"
                :key="tab.id"
                @click="activeTab = tab.id"
                :class="[
                  'px-5 py-2.5 rounded-xl font-medium transition-all duration-200',
                  activeTab === tab.id
                    ? 'bg-blue-500 text-white shadow-md'
                    : 'text-gray-600 hover:bg-gray-100'
                ]"
              >
                <span class="mr-2">{{ tab.icon }}</span>
                {{ tab.label }}
              </button>
              
              <div class="w-px h-8 bg-gray-200 mx-2"></div>
              
              <div class="flex items-center mr-4">
                <span class="text-gray-500 mr-2">👤</span>
                <span class="text-gray-700 font-medium">{{ store.currentUsername }}</span>
              </div>
              
              <button
                @click="logout"
                class="px-5 py-2.5 text-red-500 hover:bg-red-50 rounded-xl font-medium transition-colors"
              >
                退出登录
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
              <div class="space-y-2">
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
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <transition name="fade" mode="out-in">
          <DeviceList v-if="activeTab === 'devices'" key="devices" />
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
import DeviceList from './components/DeviceList.vue'
import AddDevice from './components/AddDevice.vue'
import DeviceTypeManager from './components/DeviceTypeManager.vue'
import AccountSettings from './components/AccountSettings.vue'
import DeviceDetailModal from './components/DeviceDetailModal.vue'

const activeTab = ref('devices')
const mobileMenuOpen = ref(false)

const tabs = computed(() => [
  { id: 'devices', label: '设备列表', icon: '📱' },
  { id: 'add', label: '添加设备', icon: '➕' },
  { id: 'types', label: '设备类型', icon: '📦' },
  { id: 'settings', label: '账户设置', icon: '⚙️' },
])

const logout = () => {
  if (confirm('确定要退出登录吗？')) {
    store.logout()
    activeTab.value = 'devices'
  }
}
</script>

<style>
/* 过渡动画 */
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
