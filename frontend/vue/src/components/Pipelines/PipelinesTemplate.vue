<template>
  <slot name="items" />
  <div class="flex justify-evenly items-center mt-10 h-[600px]">
    <div class="pipeline-column-block">
      <span class="pipeline-column-header">{{ headers[0] }}</span>
      <slot name="first-column" />
    </div>
    <Divider layout="vertical" />
    <div class="pipeline-column-block">
      <span class="pipeline-column-header">{{ headers[1] }}</span>
      <slot name="second-column" />
    </div>
    <Divider layout="vertical" />
    <div class="pipeline-column-block">
      <span class="pipeline-column-header">{{ headers[2] }}</span>
      <slot name="third-column" />
    </div>
  </div>
</template>

<script setup lang="ts">
import Divider from 'primevue/divider'

interface Props {
  headers?: string[];
}

defineProps<Props>()
</script>

<style scoped>

</style>