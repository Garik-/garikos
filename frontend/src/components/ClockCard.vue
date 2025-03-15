<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import timeImage from '@/assets/time.png'
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
  <div class="card d-flex flex-column">
    <div class="row row-0 flex-fill align-items-center">
      <div class="col-md-3 p-3">
        <img :src="timeImage" class="w-100 h-100 object-cover" alt="Card side image" />
      </div>
      <div class="col">
        <h3 class="card-title h1 mb-1" style="font-size: 1.85rem">{{ time }}</h3>
        <div>{{ date }}</div>
      </div>
    </div>
  </div>
</template>
