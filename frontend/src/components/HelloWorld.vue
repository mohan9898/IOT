<script setup>
defineProps({
  msg: String,
})
</script>

<template>
  <div class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow duration-300">
    <h2 class="text-xl font-semibold text-gray-800 mb-3">{{ msg }}</h2>
    <p class="text-gray-600 mb-4">点击下方按钮进行操作</p>
    <div class="space-y-2">
      <button class="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition-colors">
        查看详情
      </button>
    </div>
  </div>
</template>

<style scoped>
</style>