<script setup lang="ts">
import { computed } from 'vue'
import { formatter } from '@/utils/formatter'
import { FULL_CHARGE_VOLTAGE } from '@/config/constants'

const props = defineProps<{
  value: number
}>()

const percent = computed(() => {
  return (props.value / FULL_CHARGE_VOLTAGE) * 100
})

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
