<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="not-found">
    <div class="error-code">404</div>
    <h1 class="error-title">Страница не найдена</h1>
    <p class="error-message">Запрашиваемая страница не существует</p>
    
    <button @click="goHome" class="btn btn-primary">
      Вернуться на главную
    </button>
  </div>
</template>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.not-found {
  padding: $spacing-xl 0;
  text-align: center;
  min-height: 50vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.error-code {
  font-family: $mono-font;
  font-size: 72px;
  color: $accent-color;
  font-weight: bold;
  margin-bottom: $spacing-md;
  line-height: 1;
}

.error-title {
  font-family: $mono-font;
  font-size: $font-size-heading;
  color: $text-primary;
  margin: 0 0 $spacing-md 0;
  font-weight: normal;
}

.error-message {
  font-family: $mono-font;
  font-size: $font-size-base;
  color: $text-secondary;
  margin-bottom: $spacing-xl;
}

.btn {
  margin-top: $spacing-lg;
}
</style>
