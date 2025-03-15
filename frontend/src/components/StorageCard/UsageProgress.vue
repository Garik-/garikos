<script setup lang="ts">
import { computed } from 'vue'
import { formatBytes } from '@/utils/formatter'

const props = defineProps<{
  free: number
  usedPercent: number
  isLoading: boolean
}>()

const free = computed(() => formatBytes(props.free))
const usedPercent = computed(() => 'width: ' + props.usedPercent.toFixed(2) + '%')
</script>

<template>
  <div v-if="isLoading" class="placeholder-glow">
    <div class="placeholder col-12 placeholder-lg"></div>
  </div>

  <template v-else>
    <div class="progress progress-separated mb-3">
      <div
        class="progress-bar bg-primary"
        role="progressbar"
        :style="usedPercent"
        aria-label="Использовано"
      ></div>
    </div>
    <div class="row">
      <div class="col-auto d-flex align-items-center ps-2">
        <span class="legend me-2"></span>
        <span>Свободно</span>
        <span class="d-none d-md-inline d-lg-none d-xxl-inline ms-2 text-secondary">{{
          free
        }}</span>
      </div>
    </div>
  </template>
</template>
