<template>
  <Message
    :severity="message.status === 'success'? 'success' : 'error'"
    :pt="{
      wrapper: {
        class: '!p-2'
      }
    }"
  >
    {{ message.text }}
  </Message>
</template>

<script setup lang="ts">
import Message from 'primevue/message'
import type { DataMessage } from '@/types'

interface Props {
  message: DataMessage;
}

defineProps<Props>()
</script>

<style scoped>

</style>