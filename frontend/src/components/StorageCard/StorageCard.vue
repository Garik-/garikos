<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { DISK_ENDPOINT_URL } from '@/config/constants'
import UsageLabel from './UsageLabel.vue'
import UsageProgress from './UsageProgress.vue'

const props = defineProps<{
  path: string
}>()

const data = ref({
  isLoading: true,
  isError: false,
  total: 0,
  used: 0,
  usedPercent: 0.0,
  free: 0,
})

onMounted(async () => {
  try {
    const response = await fetch(DISK_ENDPOINT_URL + `?path=${encodeURIComponent(props.path)}`)
    if (!response.ok) {
      throw new Error('Ошибка при загрузке данных')
    }
    const { used, total, usedPercent, free } = await response.json()
    data.value = { used, total, usedPercent, free, isLoading: false, isError: false }
  } catch (err) {
    console.error(err)
    data.value.isError = true
  }
})
</script>
<template>
  <div class="card" v-if="!data.isError">
    <div class="card-header border-0 pb-0">
      <h3 class="card-title">Твердотельный накопитель</h3>
    </div>
    <div class="card-body">
      <p class="mb-0">
        Файловая система <strong>{{ props.path }}</strong>
      </p>
      <UsageLabel :isLoading="data.isLoading" :used="data.used" :total="data.total" />
      <UsageProgress
        :isLoading="data.isLoading"
        :free="data.free"
        :usedPercent="data.usedPercent"
      />
    </div>
  </div>
</template>
