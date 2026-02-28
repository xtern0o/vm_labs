<script setup>
import { useRoute } from 'vue-router'
import Breadcrumb from '@/components/Breadcrumb.vue'

const route = useRoute()
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <div class="container">
        <h1 class="app-title">labs</h1>
      </div>
    </header>
    
    <main class="app-main">
      <div class="container">
        <Breadcrumb v-if="route.path !== '/'" />
        <router-view />
      </div>
    </main>

    <footer class="app-footer">
      <div class="container">
        <div class="footer-content">
          <span class="footer-item">Карнажицкий М. Р.</span>
          <span class="footer-sep">|</span>
          <span class="footer-item">P3211</span>
          <span class="footer-sep">|</span>
          <span class="footer-item">2026</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<style lang="scss">
@import '@/styles/variables.scss';

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: $bg-secondary;
  border-bottom: $border-width solid $border-color;
  padding: $spacing-lg 0;
}

.app-title {
  font-size: $font-size-heading;
  color: $accent-color;
  font-weight: normal;
  letter-spacing: 1px;
  margin: 0;
}

.app-main {
  flex: 1;
  padding: $spacing-xl 0;
}

.app-footer {
  background: $bg-secondary;
  border-top: $border-width solid $border-color;
  padding: $spacing-lg 0;
  margin-top: auto;
}

.footer-content {
  font-family: $mono-font;
  font-size: $font-size-small;
  color: $text-secondary;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-item {
  white-space: nowrap;
}

.footer-sep {
  display: none;
}
</style>

