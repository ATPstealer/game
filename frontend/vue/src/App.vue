<template>
  <Header />
  <div class="px-0 md:px-16 py-4 relative h-full">
    <RouterView />
  </div>
  <div class="footer">
    футер
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import Header from '@/components/Header/Header.vue'
</script>

<style scoped>
.footer {
  @apply mt-auto;
}
</style>
