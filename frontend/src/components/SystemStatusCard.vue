<script setup lang="ts">
import { onMounted, useTemplateRef } from 'vue'
import ApexCharts from 'apexcharts'
import { formatBytes } from '@/utils/formatter'

let cpuChart: ApexCharts | null = null
const cpuEl = useTemplateRef('cpu-activity')

let memChart: ApexCharts | null = null
const memEl = useTemplateRef('mem-activity')

const chartOptions = {
  chart: {
    type: 'radialBar',
    fontFamily: 'inherit',
    animations: {
      enabled: false,
    },
    cpu: {
      enabled: true,
    },
  },
  tooltip: {
    enabled: false,
  },
  stroke: {
    lineCap: 'round',
  },
  plotOptions: {
    radialBar: {
      startAngle: -135,
      endAngle: 135,
      hollow: {
        margin: 0,
        size: '70%',
      },
      track: {
        startAngle: -135,
        endAngle: 135,
        margin: 0,
        background: '#dce1e7',
      },
      dataLabels: {
        show: true,
        name: {
          color: tabler.getColor('secondary'),
        },
        value: {
          offsetY: 28,
          fontWeight: '600',
          show: true,
        },
      },
    },
  },
  colors: [tabler.getColor('primary')],
}

onMounted(() => {
  cpuChart = new ApexCharts(
    cpuEl.value,
    Object.assign({}, chartOptions, {
      series: [35],
      labels: ['ЦПУ'],
    }),
  )

  memChart = new ApexCharts(
    memEl.value,
    Object.assign({}, chartOptions, {
      series: [78],
      labels: ['ЗУПВ'],
    }),
  )

  cpuChart.render()
  memChart.render()
})
</script>
<template>
  <div class="card">
    <div class="card-header border-0 pb-0">
      <h3 class="card-title">Состояние системы</h3>
    </div>
    <div class="card-body pt-0">
      <div class="row">
        <div class="col-6">
          <div ref="cpu-activity"></div>
          <h4 class="card-text text-center">26,5°</h4>
        </div>
        <div class="col-6">
          <div ref="mem-activity"></div>
          <h4 class="card-text text-center">{{ formatBytes(56.25 * 1024 * 1024) }}</h4>
        </div>
      </div>
    </div>
  </div>
</template>
