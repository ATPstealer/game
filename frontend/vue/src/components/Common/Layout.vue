<template>
  <div class="flex flex-col md:flex-row h-full">
    <Button
      v-if="!showOptions"
      class="fixed bottom-0 md:bottom-auto md:top-1/3 left-0 z-10"
      :label="t('layout.settings')"
      :pt="{
        label: {
          class: 'md:w-[1px] whitespace-pre-line break-all -ml-1.5'
        }
      }"
      @click="showOptions = true"
    />
    <div
      class="sidebar options-sidebar"
      :class="{'show-options' : showOptions}"
    >
      <Button
        v-if="showOptions"
        class="hide-button top-1/2 -right-4 hidden md:block"
        icon="pi pi-angle-double-left"
        rounded
        size="small"
        @click="showOptions = false"
      />
      <Button
        v-if="showOptions"
        class="hide-button -bottom-4 right-1/2 translate-x-1/2 md:hidden"
        icon="pi pi-angle-double-up"
        rounded
        size="small"
        @click="showOptions = false"
      />
      <p class="font-bold text-xl mb-4">
        {{ t('layout.settings') }}
      </p>
      <slot name="options" />
    </div>
    <div class="w-full">
      <slot />
    </div>
    <Button
      v-if="!showHelp"
      class="fixed bottom-0 md:bottom-auto md:top-1/3 right-0 z-10"
      :label="t('layout.help')"
      :pt="{
        label: {
          class: 'md:w-[1px] whitespace-pre-line break-all -ml-1.5'
        }
      }"
      @click="showHelp = true"
    />
    <div
      :class="showHelp ? 'sidebar help-sidebar show-help' : 'hidden'"
    >
      <Button
        v-if="showHelp"
        class="hide-button top-1/2 -left-4 hidden md:block"
        icon="pi pi-angle-double-right"
        rounded
        size="small"
        @click="showHelp = false"
      />
      <Button
        v-if="showHelp"
        class="hide-button -top-4 right-1/2 translate-x-1/2 md:hidden"
        icon="pi pi-angle-double-down"
        rounded
        size="small"
        @click="showHelp = false"
      />
      <span class="font-bold text-xl">{{ t('layout.help') }}</span>
      <div>
        <slot name="help" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

interface Props {
  showOptionsProp?: boolean;
  showHelpProp?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  showOptionsProp: true,
  showHelpProp: false
})

const showOptions = ref<boolean>(props.showOptionsProp)
const showHelp = ref<boolean>(props.showHelpProp)

const { t } = useI18n()
</script>

<style scoped>
.sidebar {
  @apply flex-col border-0 border-solid border-gray-200 p-4 pb-8 hidden
  md:flex md:pr-8 md:pb-4 relative transition-all md:w-[280px] md:min-w-[280px];
}

.options-sidebar {
  @apply md:border-r md:mr-16 md:-ml-[344px] border-b md:border-b-0 mb-8 md:mb-0;
}

.help-sidebar {
  @apply md:border-l md:ml-16 md:-mr-[344px] border-t md:border-t-0 mt-8 md:mt-0;
}

.show-options {
  @apply flex md:-ml-16 md:mr-4 ;
}

.show-help {
  @apply flex md:-mr-16 md:ml-4 ;
}

.hide-button {
  @apply absolute p-2 w-8 h-8 z-10;
}
</style>
