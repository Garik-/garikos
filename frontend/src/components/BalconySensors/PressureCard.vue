<script setup lang="ts">
import { computed, onMounted, watch, useTemplateRef } from 'vue'
import { formatter } from '@/utils/formatter'
import type { SeriesData } from '@/utils/charts'
import ApexCharts from 'apexcharts'
import { ru } from '@/utils/charts'

const props = defineProps<{
  value: number
  series: SeriesData
}>()

const formattedValue = computed(() => {
  return formatter.format(props.value)
})

const pressureSeres = (data: SeriesData) => [
  {
    name: '',
    data,
  },
]

let pressureChart: ApexCharts | null = null
const chart = useTemplateRef('chart-pressure')

onMounted(() => {
  pressureChart = new ApexCharts(chart.value, {
    chart: {
      locales: [ru],
      defaultLocale: 'ru',
      type: 'bar',
      fontFamily: 'inherit',
      height: 40,
      sparkline: {
        enabled: true,
      },
      animations: {
        enabled: false,
      },
    },
    plotOptions: {
      bar: {
        columnWidth: '50%',
      },
    },
    dataLabels: {
      enabled: false,
    },
    fill: {
      opacity: 1,
    },
    series: pressureSeres(props.series),
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

  pressureChart.render()
})

watch(
  () => props.series,
  (data) => {
    pressureChart?.updateSeries(pressureSeres(data))
  },
  { deep: true }, // deep: true, если data — это объект или массив
)
</script>

<template>
  <div class="card">
    <div class="card-body">
      <div class="subheader">Давление</div>
      <div class="d-flex align-items-baseline">
        <div class="h1 mb-0 me-2">{{ formattedValue }}</div>
        <div class="me-auto">мм рт.ст.</div>
      </div>
      <div ref="chart-pressure" class="chart-sm"></div>
    </div>
  </div>
</template>
