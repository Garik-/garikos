<script setup lang="ts">
import { computed, onMounted, useTemplateRef, watch } from 'vue'
import { formatter } from '@/utils/formatter'
import { ru } from '@/utils/charts'
import ApexCharts from 'apexcharts'
import type { SeriesData } from '@/utils/charts'

const props = defineProps<{
  value: number
  series: SeriesData
}>()

const formattedValue = computed(() => {
  return formatter.format(props.value)
})

const temperatureSeries = (data: SeriesData) => [
  {
    name: '',
    data,
  },
]

let temperatureChart: ApexCharts | null = null
const chart = useTemplateRef('chart-temperature-bg')

onMounted(() => {
  temperatureChart = new ApexCharts(chart.value, {
    chart: {
      locales: [ru],
      defaultLocale: 'ru',
      type: 'area',
      fontFamily: 'inherit',
      height: 40,
      sparkline: {
        enabled: true,
      },
      animations: {
        enabled: false,
      },
    },
    dataLabels: {
      enabled: false,
    },
    fill: {
      opacity: 0.16,
      type: 'solid',
    },
    stroke: {
      width: 2,
      lineCap: 'round',
      curve: 'smooth',
    },
    series: temperatureSeries(props.series),
    tooltip: {
      theme: 'dark',
      x: {
        format: 'dd.MM.yyyy, HH:mm',
      },
    },
    grid: {
      strokeDashArray: 4,
    },
    xaxis: {
      labels: {
        padding: 0,
      },
      tooltip: {
        enabled: false,
      },
      axisBorder: {
        show: false,
      },
      type: 'datetime',
    },
    yaxis: {
      labels: {
        padding: 4,
      },
    },
    colors: [tabler.getColor('primary')],
    legend: {
      show: false,
    },
  })

  temperatureChart.render()
})

watch(
  () => props.series,
  (data) => {
    temperatureChart?.updateSeries(temperatureSeries(data))
  },
  { deep: true }, // deep: true, если data — это объект или массив
)
</script>

<template>
  <div class="card">
    <div class="card-body">
      <div class="subheader">Температура</div>
      <div class="h1">{{ formattedValue }}°</div>
    </div>
    <div ref="chart-temperature-bg" class="chart-sm"></div>
  </div>
</template>
