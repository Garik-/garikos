<script setup lang="ts">
import IconAdguard from './icons/IconAdguard.vue'
import IconNavidrome from './icons/IconNavidrome.vue'
import IconTransmission from './icons/IconTransmission.vue'

const items = [
  {
    title: 'Navidrome',
    url: 'https://music.raspberrypi.local/',
    icon: IconNavidrome,
  },
  {
    title: 'Transmission',
    url: 'http://raspberrypi.local:9091/',
    icon: IconTransmission,
  },
  {
    title: 'AdGuard Home',
    url: 'http://raspberrypi.local:18888/',
    icon: IconAdguard,
  },
  /*{
    title: 'Grafana',
    url: 'http://raspberrypi.local:3000/',
    icon: IconGrafana,
  },
  {
    title: 'Victoria Metrics',
    url: 'http://raspberrypi.local:8428/vmui/',
    icon: IconVictoriaMetrics,
  },*/
]
</script>
<template>
  <div class="page-wrapper">
    <div class="page-header d-print-none">
      <div class="container-xl">
        <div class="page-title">Прикладные программы</div>
      </div>
    </div>

    <div class="page-body">
      <div class="container-xl">
        <div class="row row-deck row-cards">
          <div v-for="(item, index) in items" :key="index" class="col-3 col-sm-3 col-md-3 col-lg-3">
            <div class="card h-100">
              <div
                class="card-body d-flex align-items-center justify-content-center p-2 p-md-4 p-lg-5"
              >
                <a
                  :href="item.url"
                  target="_blank"
                  rel="nofollow noopener"
                  :title="item.title"
                  class="app-link"
                >
                  <component :is="item.icon" />
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.app-link {
  display: block;
  width: 100%;
}

.app-link svg {
  width: 100%;
  height: auto;
}

.app-card .card-body {
  padding: 0.5rem;
  /* минимальный padding */
  width: 100%;
  aspect-ratio: 1;
  /* квадрат */
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
