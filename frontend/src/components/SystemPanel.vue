<script setup lang="ts">
import SystemStatusCard from '@/components/SystemStatusCard.vue'
import ProcessesCard from '@/components/ProcessesCard.vue'
import StorageCard from '@/components/StorageCard/StorageCard.vue'
import ClockCard from '@/components/ClockCard.vue'
</script>

<template>
  <div class="row row-deck row-cards">
    <div class="col-12">
      <ClockCard />
    </div>
    <div class="col-12">
      <SystemStatusCard />
    </div>
    <div class="col-12">
      <StorageCard path="/" />
    </div>
    <div class="col-12">
      <StorageCard path="/mnt/music" />
    </div>
    <div class="col-12">
      <ProcessesCard :filters="['istatus', 'temperature-sensor', 'garikos', 'AdGuardHome']" />
    </div>
  </div>
</template>
