<script setup lang="ts">
import { computed } from 'vue'
import { formatter } from '@/utils/formatter'
import { MIN_CHARGE_VOLTAGE, MAX_CHARGE_VOLTAGE } from '@/config/constants'
import type { SeriesData } from '@/utils/charts'

const props = defineProps<{
  value: number
  series: SeriesData
}>()

type DayWindow = 'morning' | 'day' | 'evening'
type Bucket = { sum: number; count: number }
type DayBuckets = Record<DayWindow, Bucket>

function getWindow(date: Date) {
  const hour = date.getHours()

  if (hour >= 6 && hour < 14) return 'morning'
  if (hour >= 14 && hour < 19) return 'day'
  return 'evening' // 19–6
}

function aggregateByDayAndWindow(seriesData: SeriesData) {
  const days = new Map<string, DayBuckets>()

  for (const [ts, value] of seriesData) {
    const d = new Date(ts)

    // ключ дня: YYYY-MM-DD (локальное время)
    const dayKey = d.toISOString().slice(0, 10)

    if (!days.has(dayKey)) {
      days.set(dayKey, {
        morning: { sum: 0, count: 0 },
        day: { sum: 0, count: 0 },
        evening: { sum: 0, count: 0 },
      })
    }

    const window = getWindow(d)
    const dayBuckets = days.get(dayKey)
    if (!dayBuckets) continue
    const bucket = dayBuckets[window]

    bucket.sum += value
    bucket.count++
  }

  return days
}

function avg(sum: number, count: number) {
  return count > 0 ? sum / count : null
}

function calculateDailyConsumption(seriesData: SeriesData) {
  const days = aggregateByDayAndWindow(seriesData)
  const sortedDays = [...days.keys()].sort()

  let totalConsumption = 0
  let daysCount = 0

  for (let i = 0; i < sortedDays.length - 1; i++) {
    const todayKey = sortedDays[i]
    const nextDayKey = sortedDays[i + 1]
    if (!todayKey || !nextDayKey) continue

    const today = days.get(todayKey)
    const nextDay = days.get(nextDayKey)
    if (!today || !nextDay) continue

    const morning = avg(today.morning.sum, today.morning.count)
    const day = avg(today.day.sum, today.day.count)
    const evening = avg(today.evening.sum, today.evening.count)
    const nextMorning = avg(nextDay.morning.sum, nextDay.morning.count)

    if (morning === null || day === null || evening === null || nextMorning === null) {
      continue
    }

    let consumption = 0

    if (day < morning) consumption += morning - day
    if (evening < day) consumption += day - evening
    if (nextMorning < evening) consumption += evening - nextMorning

    totalConsumption += consumption
    daysCount++
  }

  return daysCount > 0 ? totalConsumption / daysCount : null
}

function estimateDaysLeft(
  currentMv: number,
  dailyConsumptionMv: number | null,
  minMv = MIN_CHARGE_VOLTAGE,
) {
  if (!dailyConsumptionMv || dailyConsumptionMv <= 0) return null
  return (currentMv - minMv) / dailyConsumptionMv
}

const getBatteryPercent = (voltage_mV: number) => {
  if (voltage_mV >= MAX_CHARGE_VOLTAGE) return 100
  if (voltage_mV <= MIN_CHARGE_VOLTAGE) return 0

  // LUT в милливольтах (отсортирован по убыванию!)
  const lut_mV: readonly number[] = [
    MAX_CHARGE_VOLTAGE,
    4100,
    4000,
    3900,
    3800,
    3700,
    3600,
    3500,
    3300,
    3000,
    MIN_CHARGE_VOLTAGE,
  ]
  const lut_pct: readonly number[] = [100, 95, 85, 75, 60, 40, 20, 10, 5, 1, 0]

  // Бинарный поиск
  let low = 0
  let high = lut_mV.length - 1
  while (low <= high) {
    const mid = (low + high) >>> 1
    const midVal = lut_mV[mid]
    if (midVal === undefined) break
    if (midVal === voltage_mV) return lut_pct[mid] ?? 0
    if (midVal < voltage_mV) high = mid - 1
    else low = mid + 1
  }

  // Интерполяция между high и low (теперь high < low)
  const i = Math.max(0, Math.min(high, lut_mV.length - 2))
  const upperV = lut_mV[i]
  const lowerV = lut_mV[i + 1]
  const upperP = lut_pct[i]
  const lowerP = lut_pct[i + 1]
  if (upperV === undefined || lowerV === undefined || upperP === undefined || lowerP === undefined) {
    return 0
  }

  const rangeV = upperV - lowerV
  const rangeP = upperP - lowerP
  if (rangeV === 0) return lowerP

  return lowerP + Math.round(((voltage_mV - lowerV) * rangeP) / rangeV)
}

const percent = computed<number>(() => getBatteryPercent(props.value))

const dailyConsumptionMv = computed(() => calculateDailyConsumption(props.series))

const averageLabel = computed(() => {
  if (dailyConsumptionMv.value !== null) {
    return formatter.format(dailyConsumptionMv.value)
  } else {
    return '0'
  }
})

const daysLeft = computed(() => estimateDaysLeft(props.value, dailyConsumptionMv.value))
const daysLeftFormatted = computed(() =>
  daysLeft.value !== null ? formatter.format(daysLeft.value) : '0',
)

const label = computed(() => {
  return formatter.format(props.value)
})

const objectOfAttrs = computed(() => ({
  style: 'width:' + percent.value.toFixed(2) + '%',
  'aria-valuenow': +percent.value.toFixed(2),
  'aria-label': label.value,
}))
</script>

<!--
<template>
  <div class="card">
    <div class="card-body d-flex flex-column">
      <div class="subheader">Батарея</div>
      <div class="d-flex align-items-baseline">
        <div class="h1 mb-0 me-2">{{ label }}</div>
        <div class="me-auto">мВ</div>
      </div>
      <div class="text-secondary">{{ averageLabel }} мВ/сут</div>
      <div class="text-secondary">~{{ daysLeft }} дней</div>

      <div class="progress progress-sm mt-auto">
        <div v-bind="objectOfAttrs" class="progress-bar bg-primary" role="progressbar" aria-valuemin="0"
          aria-valuemax="100">
          <span class="visually-hidden">{{ label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
-->

<template>
  <div class="card">
    <div class="card-body d-flex flex-column">
      <!-- Заголовок -->
      <div class="subheader text-muted">Батарея</div>

      <!-- Основное значение -->
      <div class="d-flex align-items-baseline">
        <div class="h1 mb-0 me-2">{{ label }}</div>
        <div class="me-auto">мВ</div>
      </div>

      <!-- Метрики -->
      <div v-if="dailyConsumptionMv !== null" class="text-muted small">
        {{ averageLabel }} мВ / сутки
      </div>
      <div v-if="daysLeft !== null" class="text-muted small">≈ {{ daysLeftFormatted }} дней</div>

      <!-- Прогресс -->
      <div class="mt-2">
        <div class="progress progress-sm">
          <div
            v-bind="objectOfAttrs"
            class="progress-bar bg-primary"
            role="progressbar"
            aria-valuemin="0"
            aria-valuemax="100"
          />
        </div>
      </div>
    </div>
  </div>
</template>
