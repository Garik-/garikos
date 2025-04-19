<script setup lang="ts">
import { computed } from 'vue'
import { formatter } from '@/utils/formatter'
import { MIN_CHARGE_VOLTAGE, MAX_CHARGE_VOLTAGE } from '@/config/constants'

const props = defineProps<{
  value: number
}>()

const getBatteryPercent = (voltage_mV: number) => {
  if (voltage_mV >= 4200) return 100
  if (voltage_mV <= 2500) return 0

  // LUT в милливольтах (отсортирован по убыванию!)
  const lut_mV = [4200, 4100, 4000, 3900, 3800, 3700, 3600, 3500, 3300, 3000, 2500]
  const lut_pct = [100, 95, 85, 75, 60, 40, 20, 10, 5, 1, 0]

  // Бинарный поиск
  let low = 0
  let high = lut_mV.length - 1
  while (low <= high) {
    const mid = (low + high) >>> 1
    const midVal = lut_mV[mid]
    if (midVal === voltage_mV) return lut_pct[mid]
    if (midVal < voltage_mV) high = mid - 1
    else low = mid + 1
  }

  // Интерполяция между high и low (теперь high < low)
  const i = high
  const rangeV = lut_mV[i] - lut_mV[i + 1]
  const rangeP = lut_pct[i] - lut_pct[i + 1]
  return lut_pct[i + 1] + Math.round(((voltage_mV - lut_mV[i + 1]) * rangeP) / rangeV)
}

const percent = computed(() => getBatteryPercent(props.value))

const label = computed(() => {
  return formatter.format(props.value)
})

const objectOfAttrs = computed(() => ({
  style: 'width:' + percent.value.toFixed(2) + '%',
  'aria-valuenow': +percent.value.toFixed(2),
  'aria-label': label.value,
}))
</script>

<template>
  <div class="card">
    <div class="card-body d-flex flex-column">
      <div class="subheader">Батарея</div>
      <div class="d-flex align-items-baseline">
        <div class="h1 mb-0 me-2">{{ label }}</div>
        <div class="me-auto">мВ</div>
      </div>

      <div class="progress progress-sm mt-auto">
        <div
          v-bind="objectOfAttrs"
          class="progress-bar bg-primary"
          role="progressbar"
          aria-valuemin="0"
          aria-valuemax="100"
        >
          <span class="visually-hidden">{{ label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
