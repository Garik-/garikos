<script setup lang="ts">
import BalconySensors from '@/components/BalconySensors/BalconySensors.vue'
import AppsPanel from './components/AppsPanel.vue'
import SystemPanel from './components/SystemPanel.vue'
</script>

<template>
  <div class="page">
    <div class="page-wrapper">
      <div class="page-body">
        <div class="container-xl">
          <!-- mobile -->
          <div class="mobile-carousel d-lg-none">
            <div class="carousel-track">
              <section class="carousel-slide">
                <SystemPanel />
              </section>

              <section class="carousel-slide">
                <AppsPanel />
              </section>

              <section class="carousel-slide">
                <BalconySensors />
              </section>
            </div>
          </div>
          <!--/mobile-->

          <!--desktop-->
          <div class="d-none d-lg-block">
            <div class="row">
              <div class="col-sm-6 col-lg-4">
                <div class="page-body">
                  <SystemPanel />
                </div>
              </div>
              <div class="col-sm-6 col-lg-8">
                <BalconySensors />
                <AppsPanel />
              </div>
            </div>
          </div>

          <!--/desktop-->
        </div>
      </div>
    </div>
  </div>
</template>
