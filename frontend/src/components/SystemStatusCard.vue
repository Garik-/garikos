<script setup lang="ts">
import { onMounted, onUnmounted, useTemplateRef, ref } from 'vue'
import ApexCharts from 'apexcharts'
import { formatBytes } from '@/utils/formatter'
import { formatter } from '@/utils/formatter'
import { createSSEConnection } from '@/services/sseService'
import { SYSTEM_SSE_URL } from '@/config/constants'

const cpuTemperature = ref('')
const memBytes = ref('')

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

let eventSource: EventSource | null = null

type Data = {
  cpu: [number]
  sensors: [
    {
      temperature: number
    },
  ]
  mem: {
    used: number
    usedPercent: number
  }
}

onMounted(() => {
  cpuChart = new ApexCharts(
    cpuEl.value,
    Object.assign({}, chartOptions, {
      series: [],
      labels: ['ЦПУ'],
    }),
  )

  memChart = new ApexCharts(
    memEl.value,
    Object.assign({}, chartOptions, {
      series: [],
      labels: ['ЗУПВ'],
    }),
  )

  cpuChart.render()
  memChart.render()

  eventSource = createSSEConnection(SYSTEM_SSE_URL, (d) => {
    console.log(d)

    cpuTemperature.value = formatter.format((d as Data).sensors[0].temperature)
    cpuChart?.updateSeries((d as Data).cpu.map((value) => Math.round(value * 100) / 100))

    memChart?.updateSeries([Math.round((d as Data).mem.usedPercent * 100) / 100])
    memBytes.value = formatBytes((d as Data).mem.used)
  })
})

onUnmounted(() => {
  eventSource?.close()
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
          <h4 class="card-text text-center">{{ cpuTemperature }}°</h4>
        </div>
        <div class="col-6">
          <div ref="mem-activity"></div>
          <h4 class="card-text text-center">{{ memBytes }}</h4>
        </div>
      </div>
    </div>
  </div>
</template>
