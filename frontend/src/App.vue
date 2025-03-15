<script setup lang="ts">
import BalconySensors from '@/components/BalconySensors/BalconySensors.vue'
import SystemStatusCard from '@/components/SystemStatusCard.vue'
import ProcessesCard from './components/ProcessesCard.vue'
import AppsPanel from './components/AppsPanel.vue'
import StorageCard from './components/StorageCard.vue'
import ClockCard from './components/ClockCard.vue'
</script>

<template>
  <div class="page">
    <div class="page-wrapper">
      <div class="page-body">
        <div class="container-xl">
          <div class="row">
            <div class="col-sm-6 col-lg-4">
              <div class="page-body">
                <div class="row row-deck row-cards">
                  <div class="col-12">
                    <ClockCard />
                  </div>
                  <div class="col-12">
                    <SystemStatusCard />
                  </div>
                  <div class="col-12">
                    <StorageCard />
                  </div>
                  <div class="col-12">
                    <ProcessesCard />
                  </div>
                </div>
              </div>
            </div>
            <div class="col-sm-6 col-lg-8">
              <BalconySensors />
              <AppsPanel />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
