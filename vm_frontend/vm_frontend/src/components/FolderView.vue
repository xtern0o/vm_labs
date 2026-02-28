<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { structure } from '@/config/pages.js'
import NotFound from '@/components/NotFound.vue'

const route = useRoute()
const router = useRouter()

// Найти текущую папку по пути
const currentFolder = computed(() => {
  const pathParts = route.path.split('/').filter(Boolean)
  
  if (pathParts.length === 0) {
    return null // Корневая директория
  }
  
  let current = { children: structure }
  
  for (const part of pathParts) {
    const found = current.children?.find(item => item.id === part)
    if (found && found.type === 'folder') {
      current = found
    } else {
      return undefined // Путь не найден
    }
  }
  
  return current
})

const items = computed(() => {
  return currentFolder.value ? currentFolder.value.children : structure
})

// Проверка существования пути
const pathExists = computed(() => {
  return currentFolder.value !== undefined
})

const openItem = (item) => {
  router.push(item.path)
}
</script>

<template>
  <NotFound v-if="!pathExists" />
  <div v-else class="folder-view">
    <h1 class="folder-title">
      {{ currentFolder ? currentFolder.name : 'Index of /' }}
    </h1>
    
    <div class="items-list">
      <div class="list-header">
        <span class="col-name">Name</span>
        <span class="col-desc">Description</span>
      </div>
      <div
        v-for="item in items"
        :key="item.id"
        class="item-row"
        @click="openItem(item)"
      >
        <span class="item-name">
          {{ item.type === 'folder' ? '[DIR]' : '[PAGE]' }} {{ item.id }}{{ item.type === 'folder' ? '/' : '' }}
        </span>
        <span class="item-desc">{{ item.name }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.folder-view {
  padding: $spacing-md 0;
}

.folder-title {
  font-family: $mono-font;
  font-size: $font-size-heading;
  color: $text-primary;
  margin: 0 0 $spacing-xl 0;
  font-weight: normal;
  border-bottom: $border-width solid $border-color;
  padding-bottom: $spacing-md;
}

.items-list {
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

.item-row {
  display: flex;
  padding: $spacing-sm 0;
  cursor: pointer;
  transition: background-color $transition-fast;
  border-bottom: $border-width solid $border-color;

  &:hover {
    background-color: rgba($accent-color, 0.05);
    
    .item-name {
      color: $accent-color;
    }
  }
}

.item-name {
  flex: 0 0 250px;
  color: $accent-color;
  font-weight: normal;
  transition: color $transition-fast;
}

.item-desc {
  flex: 1;
  color: $text-secondary;
  font-size: $font-size-small;
}
</style>
