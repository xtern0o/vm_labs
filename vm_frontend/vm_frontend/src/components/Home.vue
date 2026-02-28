<script setup>
import { useRouter } from 'vue-router'
import { structure } from '@/config/pages.js'

const router = useRouter()

const openItem = (item) => {
  router.push(item.path)
}
</script>

<template>
  <div class="home-container">
    <h1 class="home-title">Index of /</h1>
    
    <div class="labs-list">
      <div class="list-header">
        <span class="col-name">Name</span>
        <span class="col-desc">Description</span>
      </div>
      <div
        v-for="item in structure"
        :key="item.id"
        class="lab-item"
        @click="openItem(item)"
      >
        <span class="lab-name">
          {{ item.type === 'folder' ? '[DIR]' : '[PAGE]' }} {{ item.id }}{{ item.type === 'folder' ? '/' : '' }}
        </span>
        <span class="lab-desc">{{ item.name }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.home-container {
  padding: $spacing-md 0;
}

.home-title {
  font-family: $mono-font;
  font-size: $font-size-heading;
  color: $text-primary;
  margin: 0 0 $spacing-xl 0;
  font-weight: normal;
  border-bottom: $border-width solid $border-color;
  padding-bottom: $spacing-md;
}

.labs-list {
  font-family: $mono-font;
}

.list-header {
  display: flex;
  padding: $spacing-sm 0;
  border-bottom: $border-width solid $border-color;
  color: $text-secondary;
  font-size: $font-size-small;
  text-transform: uppercase;
  
  .col-name {
    flex: 0 0 250px;
  }
  
  .col-desc {
    flex: 1;
  }
}

.lab-item {
  display: flex;
  padding: $spacing-sm 0;
  cursor: pointer;
  transition: background-color $transition-fast;
  border-bottom: $border-width solid $border-color;

  &:hover {
    background-color: rgba($accent-color, 0.05);
    
    .lab-name {
      color: $accent-color;
    }
  }
}

.lab-name {
  flex: 0 0 250px;
  color: $accent-color;
  font-weight: normal;
  transition: color $transition-fast;
}

.lab-desc {
  flex: 1;
  color: $text-secondary;
  font-size: $font-size-small;
}
</style>
