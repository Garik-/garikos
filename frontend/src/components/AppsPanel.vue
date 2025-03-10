<script setup lang="ts">
import IconAdguard from './icons/IconAdguard.vue'
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
          <div class="col-sm-3 col-lg-3">
            <div class="card">
              <div class="card-body">
                <a href="http://raspberrypi.local:18888/" target="_blank" class="app-link">
                  <IconAdguard />
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
.app-link svg {
  width: 100%;
  height: auto;
}
</style>
