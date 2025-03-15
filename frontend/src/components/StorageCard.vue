<template>
  <div class="card">
    <div class="card-header border-0 pb-0">
      <h3 class="card-title">Твердотельный накопитель</h3>
    </div>
    <div class="card-body">
      <p class="mb-3">
        Используется <strong>6,854 Мбайт</strong> из 8 Гбайт доступного пространства
      </p>
      <div class="progress progress-separated mb-3">
        <div
          class="progress-bar bg-primary"
          role="progressbar"
          style="width: 44%"
          aria-label="Regular"
        ></div>
        <div
          class="progress-bar bg-info"
          role="progressbar"
          style="width: 19%"
          aria-label="System"
        ></div>
        <div
          class="progress-bar bg-success"
          role="progressbar"
          style="width: 9%"
          aria-label="Shared"
        ></div>
      </div>
      <div class="row">
        <div class="col-auto d-flex align-items-center pe-2">
          <span class="legend me-2 bg-primary"></span>
          <span>Regular</span>
          <span class="d-none d-md-inline d-lg-none d-xxl-inline ms-2 text-secondary">915MB</span>
        </div>
        <div class="col-auto d-flex align-items-center px-2">
          <span class="legend me-2 bg-info"></span>
          <span>System</span>
          <span class="d-none d-md-inline d-lg-none d-xxl-inline ms-2 text-secondary">415MB</span>
        </div>
        <div class="col-auto d-flex align-items-center px-2">
          <span class="legend me-2 bg-success"></span>
          <span>Shared</span>
          <span class="d-none d-md-inline d-lg-none d-xxl-inline ms-2 text-secondary">201MB</span>
        </div>
        <div class="col-auto d-flex align-items-center ps-2">
          <span class="legend me-2"></span>
          <span>Free</span>
          <span class="d-none d-md-inline d-lg-none d-xxl-inline ms-2 text-secondary">612MB</span>
        </div>
      </div>
    </div>
  </div>
</template>
