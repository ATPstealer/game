<template>
  <Header />
  <div class="px-0 md:px-16 py-4 relative h-full">
    <ConfirmPopup
      :pt="{
        message: {
          class: 'max-w-[250px]'
        }
      }"
    />
    <RouterView />
    <Toast />
  </div>
  <div class="footer">
    футер
  </div>
</template>

<script setup lang="ts">
import ConfirmPopup from 'primevue/confirmpopup'
import Toast from 'primevue/toast'
import { RouterView } from 'vue-router'
import Header from '@/components/Header/Header.vue'
</script>

<style scoped>
.footer {
  @apply mt-auto;
}
</style>
