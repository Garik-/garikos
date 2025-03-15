<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
const time = ref('')
const date = ref('')

const updateTime = () => {
  const now = new Date()
  time.value = now.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
  date.value = now.toLocaleDateString('ru-RU', { weekday: 'long', day: 'numeric', month: 'long' })
}

const interval = setInterval(updateTime, 1_000)
onMounted(updateTime)
onUnmounted(() => clearInterval(interval))
</script>
<template>
  <div class="card">
    <div class="card-header border-0 pb-0">
      <h3 class="card-title h1 mb-2" style="font-size: 1.85rem">{{ time }}</h3>
    </div>
    <div class="card-body pt-0">{{ date }}</div>
  </div>
</template>
