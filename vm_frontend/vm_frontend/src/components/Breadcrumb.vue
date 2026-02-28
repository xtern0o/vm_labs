<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const breadcrumbs = computed(() => {
  const crumbs = [{ name: '/', path: '/' }]
  
  const currentPath = route.path
  if (currentPath !== '/') {
    const parts = currentPath.split('/').filter(Boolean)
    let pathSoFar = ''
    
    parts.forEach(part => {
      pathSoFar += `/${part}`
      crumbs.push({ name: part, path: pathSoFar })
    })
  }
  
  return crumbs
})

const navigate = (path) => {
  router.push(path)
}
</script>

<template>
  <nav class="breadcrumb">
    <span v-for="(crumb, index) in breadcrumbs" :key="index" class="breadcrumb-item">
      <a v-if="crumb.path !== route.path" @click="navigate(crumb.path)" class="breadcrumb-link">
        {{ crumb.name }}
      </a>
      <span v-else class="breadcrumb-current">{{ crumb.name }}</span>
      <span v-if="index < breadcrumbs.length - 1" class="breadcrumb-sep"> > </span>
    </span>
  </nav>
</template>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.breadcrumb {
  padding: $spacing-md 0;
  font-family: $mono-font;
  font-size: $font-size-small;
  color: $text-secondary;
  border-bottom: $border-width solid $border-color;
  margin-bottom: $spacing-lg;
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.breadcrumb-item {
  display: inline-flex;
  align-items: center;
}

.breadcrumb-link {
  color: $accent-color;
  cursor: pointer;
  text-decoration: none;
  transition: all $transition-fast;
  padding: $spacing-xs $spacing-md;
  border: $border-width solid $border-color;
  background: $bg-secondary;
  display: inline-block;
  
  &:hover {
    border-color: $accent-color;
    background: rgba($accent-color, 0.05);
  }
  
  &:active {
    background: rgba($accent-color, 0.1);
  }
}

.breadcrumb-current {
  color: $text-primary;
  padding: $spacing-xs $spacing-md;
  border: $border-width solid $border-color;
  background: $bg-primary;
  display: inline-block;
}

.breadcrumb-sep {
  margin: 0 $spacing-xs;
  color: $text-secondary;
}
</style>
