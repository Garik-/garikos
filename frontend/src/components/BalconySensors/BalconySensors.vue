<script setup lang="ts">
import TemperatureCard from './TemperatureCard.vue'
import HumidityCard from './HumidityCard.vue'
import PressureCard from './PressureCard.vue'
import VoltageCard from './VoltageCard.vue'

import { BALCONY_SSE_URL } from '@/config/constants'
import { dateToLocaleString } from '@/utils/formatter'

import { onMounted, onUnmounted, ref, computed } from 'vue'
import { createSSEConnection } from '@/services/sseService'

import type { SeriesData } from '@/utils/charts'
import type { Ref } from 'vue'

type Data = {
  current: {
    temperature: number
    humidity: number
    pressure: number
    timestamp: string
    voltage: number
  }
  chart: {
    temperature: SeriesData
    pressure: SeriesData
  }
}

const data: Ref<Data> = ref({
  current: {
    temperature: 0,
    humidity: 0,
    pressure: 0,
    timestamp: '0001-01-01T00:00:00Z',
    voltage: 0,
  },
  chart: { temperature: [], pressure: [] },
})

const currentTime = computed(() => dateToLocaleString(new Date(data.value.current.timestamp)))

let eventSource: EventSource | null = null

onMounted(() => {
  eventSource = createSSEConnection(BALCONY_SSE_URL, (d) => {
    data.value = d as Data
  })
})

onUnmounted(() => {
  eventSource?.close()
})
</script>
<template>
  <div class="page-wrapper">
    <div class="page-header d-print-none">
      <div class="container-xl">
        <div class="row g-2 align-items-center">
          <div class="col">
            <div class="page-pretitle">Сенсоры</div>
            <div class="page-title">На балконе</div>
          </div>
          <div class="col-auto text-end">
            <div class="text-secondary fs-5">последнее обновление</div>
            <div>{{ currentTime }}</div>
          </div>
        </div>
      </div>
    </div>
    <div class="page-body">
      <div class="container-xl">
        <div class="row row-deck row-cards">
          <div class="col-sm-6 col-lg-3">
            <TemperatureCard :value="data.current.temperature" :series="data.chart.temperature" />
          </div>
          <div class="col-sm-6 col-lg-3">
            <HumidityCard :value="data.current.humidity" />
          </div>
          <div class="col-sm-6 col-lg-3">
            <PressureCard :value="data.current.pressure" :series="data.chart.pressure" />
          </div>
          <div class="col-sm-6 col-lg-3">
            <VoltageCard :value="data.current.voltage" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
