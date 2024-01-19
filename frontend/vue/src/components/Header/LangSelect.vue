<template>
  <Dropdown
    v-model="lang"
    :options="languages"
    option-label="label"
    option-value="key"
    @change="changeLang"
    class="lang"
    :pt="{
      input: {
        class: 'text-white py-1'
      },
      trigger: {
        class: 'text-white py-1'
      },
      item: {
        class: 'p-2'
      },
      panel: {
        class: 'rounded-none'
      }
    }"
  />
</template>

<script setup lang="ts">
import Dropdown, { DropdownChangeEvent } from 'primevue/dropdown'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { languages } from '@/components/Header/constants'

const lang = ref<string>(localStorage.getItem('locale') || 'ru')

const { locale } = useI18n()
const changeLang = (event: DropdownChangeEvent) => {
  locale.value = event.value
  localStorage.setItem('locale', event.value)
}
</script>

<style scoped>
.lang {
  @apply border-b border-solid border-b-white rounded-none bg-gray-500 border-gray-500 items-baseline relative z-[9999];
}

.p-dropdown:not(.p-disabled).p-focus {
  @apply shadow-none;
}
</style>