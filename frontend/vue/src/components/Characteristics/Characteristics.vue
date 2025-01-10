<template>
  <Card>
    <template #content>
      <h2 class="font-bold mb-2 text-2xl">
        {{ t('user.characteristics') }}
      </h2>
      <div v-if="!isEmpty(characteristics)" class="grid gap-1">
        <p v-for="char in Object.keys(characteristics)" :key="char">
          {{ t(`user.${char}`) }}: {{ characteristics[char as keyof Characteristics] }}
        </p>
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import isEmpty from 'lodash/isEmpty'
import { storeToRefs } from 'pinia'
import Card from 'primevue/card'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Characteristics } from '@/gen'
import { useUserStore } from '@/stores/userStore'

const { t } = useI18n()

const { userData } = storeToRefs(useUserStore())

const characteristics = computed(() => {
  if (!userData.value) {
    return
  }

  return userData.value.characteristics
})

</script>

<style scoped>

</style>