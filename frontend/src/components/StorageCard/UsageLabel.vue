<script setup lang="ts">
import { computed } from 'vue'
import { formatBytes } from '@/utils/formatter'
const props = defineProps<{
  used: number
  total: number
  isLoading: boolean
}>()

const used = computed(() => formatBytes(props.used))
const total = computed(() => formatBytes(props.total))
</script>
<template>
  <p v-if="isLoading" class="mb-3 placeholder-glow">
    <span class="placeholder col-12"></span>
  </p>
  <p v-else class="mb-3">
    Используется <strong>{{ used }}</strong> из {{ total }}
  </p>
</template>
