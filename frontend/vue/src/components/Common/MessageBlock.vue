<template>
  <InlineMessage :severity="message.status === 'success'? 'success' : 'error'" :closable="false">
    {{ message.text }}
  </InlineMessage>
</template>

<script setup lang="ts">
import InlineMessage from 'primevue/inlinemessage'
import type { Message } from '@/types'

interface Props {
  message: Message;
}

defineProps<Props>()
</script>

<style scoped>

</style>