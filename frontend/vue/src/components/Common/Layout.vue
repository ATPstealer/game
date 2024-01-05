<template>
  <div class="flex flex-col md:flex-row h-full">
    <Button
        :label="t('layout.settings')"
        class="fixed bottom-0 md:bottom-auto md:top-1/3 left-0 z-10"
        v-if="!showOptions"
        :pt="{
        label: {
          class: 'md:w-[1px] whitespace-pre-line break-all -ml-1.5'
        }
      }"
        @click="showOptions = true"
    />
    <div
        class="options-sidebar"
        :class="{'show' : showOptions}"
    >
      <Button
          v-if="showOptions"
          icon="pi pi-angle-double-left"
          class="hide-button top-1/2 -right-4 hidden md:block"
          rounded
          size="small"
          @click="showOptions = false"
      />
      <Button
          v-if="showOptions"
          icon="pi pi-angle-double-up"
          class="hide-button -bottom-4 right-1/2 translate-x-1/2 md:hidden"
          rounded
          size="small"
          @click="showOptions = false"
      />
      <span class="font-bold text-xl">{{ t('layout.settings') }}</span>
      <slot name="options"/>
    </div>
    <div class="w-full">
      <slot/>
    </div>
    <Button
        :label="t('layout.help')"
        class="fixed bottom-0 md:bottom-auto md:top-1/3 right-0 z-10"
        :pt="{
        label: {
          class: 'md:w-[1px] whitespace-pre-line break-all -ml-1.5'
        }
      }"
        @click="showHelp = true"
    />
    <Sidebar v-model:visible="showHelp" position="right">
      <span class="font-bold text-xl">{{ t('layout.help') }}</span>
      <div>
        <slot name="help"/>
      </div>
    </Sidebar>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Sidebar from 'primevue/sidebar'
import {ref} from 'vue'
import {useI18n} from "vue-i18n";

const showOptions = ref<boolean>(true)
const showHelp = ref<boolean>(false)

const {t} = useI18n()
</script>

<style scoped>
.options-sidebar {
  @apply flex-col gap-8 border-0 border-b border-solid border-gray-200 p-4 pb-8 hidden
  md:flex md:pr-8 md:pb-4 md:border-r md:border-b-0 md:mr-16 relative transition-all md:w-[280px] md:min-w-[280px] md:-ml-[344px];
}

.show {
  @apply flex md:-ml-16 md:mr-4 ;

}

.hide-button {
  @apply absolute p-2 w-8 h-8 z-10;
}
</style>
