<script setup>
import { defineProps, defineEmits } from 'vue'

defineProps({
  message: {
    type: String,
    required: true
  },
  type: {
    type: String,
    default: 'error',
    validator: (value) => ['error', 'warning', 'info', 'success'].includes(value)
  }
})

const emit = defineEmits(['close'])

const handleClose = () => {
  emit('close')
}
</script>

<template>
  <div :class="['alert', `alert-${type}`]" @click="handleClose">
    {{ message }}
  </div>
</template>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.alert {
  cursor: pointer;
  user-select: none;
  
  &:hover {
    opacity: 0.8;
  }
  
  &:active {
    opacity: 0.6;
  }
}
</style>
