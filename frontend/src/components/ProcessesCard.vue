<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { Ref } from 'vue'
import { createSSEConnection } from '@/services/sseService'
import { PROC_SSE_URL } from '@/config/constants'
import { formatter, formatBytes } from '@/utils/formatter'

const props = defineProps<{
  filters?: string[]
}>()

type Item = {
  pid: number
  name: string
  cpuPercent: string
  memRSS: string
}

type Response = {
  pid: number
  name: string
  cpuPercent: number
  mem: {
    rss: number
  }
}

const items: Ref<Item[]> = ref([])

function getURL(filters?: string[]) {
  if (!filters) return PROC_SSE_URL
  let url: URL

  try {
    url = new URL(PROC_SSE_URL)
  } catch {
    url = new URL(PROC_SSE_URL, window.location.origin)
  }

  const searchParams = new URLSearchParams()
  filters.forEach((value) => {
    searchParams.append('name', value)
  })
  url.search = searchParams.toString()

  return url.toString()
}

let eventSource: EventSource | null = null
onMounted(() => {
  eventSource = createSSEConnection(getURL(props.filters), (processes) => {
    items.value = (processes as Response[]).map((p): Item => {
      return {
        pid: p.pid,
        name: p.name,
        cpuPercent: formatter.format(p.cpuPercent) + '%',
        memRSS: formatBytes(p.mem.rss),
      }
    })
  })
})

onUnmounted(() => {
  eventSource?.close()
})
</script>
<template>
  <div class="card">
    <div class="card-header">
      <h3 class="card-title">Монитор процессов</h3>
    </div>
    <div class="card-table table-responsive">
      <table class="table table-vcenter">
        <thead>
          <tr>
            <th title="PID">ИП</th>
            <th>Команда</th>
            <th>ЦПУ</th>
            <th>ЗУПВ</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.pid">
            <td>{{ item.pid }}</td>
            <td>{{ item.name }}</td>
            <td>{{ item.cpuPercent }}</td>
            <td>{{ item.memRSS }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
