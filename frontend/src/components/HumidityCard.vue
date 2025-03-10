<script setup lang="ts">
import { computed } from 'vue'
import { formatter } from '@/utils/formatter'

const props = defineProps<{
  value: number
}>()

const label = computed(() => {
  return formatter.format(props.value) + '%'
})

const objectOfAttrs = computed(() => ({
  style: 'width:' + props.value.toFixed(2) + '%',
  'aria-valuenow': +props.value.toFixed(2),
  'aria-label': label.value,
}))
</script>

<template>
  <div class="card">
    <div class="card-body d-flex flex-column">
      <div class="subheader">Влажность</div>
      <div class="h1">{{ label }}</div>

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
