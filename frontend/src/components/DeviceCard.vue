<template>
  <div
    @click="handleClick"
    class="bg-white border-2 border-gray-100 rounded-xl p-4 hover:shadow-lg hover:border-blue-200 cursor-pointer transition-all duration-300 transform hover:-translate-y-1"
  >
    <div class="flex justify-between items-start mb-3">
      <span class="text-3xl">{{ icon }}</span>
      <span
        :class="[
          'px-2 py-0.5 rounded-full text-xs font-semibold whitespace-nowrap flex-shrink-0',
          device.status === 'online' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
        ]"
      >
        {{ device.status === 'online' ? '在线' : '离线' }}
      </span>
    </div>
    <h3 class="font-bold text-gray-800 text-base mb-0.5 truncate" :title="device.name">{{ device.name }}</h3>
    <p class="text-xs text-gray-500 whitespace-nowrap">{{ typeName }}</p>
    <div class="mt-2 text-xs text-gray-400 font-mono truncate" :title="device.id">{{ device.id }}</div>
    <div v-if="device.metadata?.note" class="mt-2 text-xs text-gray-500 bg-gray-50 p-2 rounded-lg truncate" :title="device.metadata.note">
      {{ device.metadata.note }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import store from '../store'

const props = defineProps({
  device: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['select'])

const icon = computed(() => store.getDeviceIcon(props.device.type))
const typeName = computed(() => store.getDeviceTypeName(props.device.type))

const handleClick = () => {
  emit('select', props.device)
}
</script>
