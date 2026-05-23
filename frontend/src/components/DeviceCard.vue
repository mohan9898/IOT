<template>
  <div
    @click="handleClick"
    class="bg-white border-2 border-gray-100 rounded-xl p-5 hover:shadow-lg hover:border-blue-200 cursor-pointer transition-all duration-300 transform hover:-translate-y-1"
  >
    <div class="flex justify-between items-start mb-4">
      <span class="text-4xl">{{ icon }}</span>
      <span
        :class="[
          'px-3 py-1 rounded-full text-xs font-semibold',
          device.status === 'online' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
        ]"
      >
        {{ device.status === 'online' ? '在线' : '离线' }}
      </span>
    </div>
    <h3 class="font-bold text-gray-800 text-lg mb-1 truncate">{{ device.name }}</h3>
    <p class="text-sm text-gray-500">{{ typeName }}</p>
    <div class="mt-3 text-xs text-gray-400 font-mono truncate">{{ device.id }}</div>
    <div v-if="device.metadata?.note" class="mt-2 text-xs text-gray-500 bg-gray-50 p-2 rounded-lg">
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
