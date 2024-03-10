<template>
  <Card>
    <template #content>
      <h2 class="font-bold mb-2 text-2xl">
        {{ t('user.characteristics') }}
      </h2>
      <div class="grid gap-1" v-if="!isFetching">
        <p v-for="char in Object.keys(characteristics)" :key="char">
          {{ t(`user.${char}`) }}: {{ characteristics[char] }}
        </p>
      </div>
      <Loading v-else />
    </template>
  </Card>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import { computed, Ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Loading from '@/components/Common/Loading.vue'
import { useGetData } from '@/composables/useGetData'
import type { Characteristics } from '@/types'

const { t } = useI18n()
const { data, isFetching } = useGetData('/user/data')

const characteristics: Ref<Characteristics> = computed(() => {
  return {
    memory: data.value?.characteristics.memory,
    intelligence: data.value?.characteristics.intelligence,
    attention: data.value?.characteristics.attention,
    wits: data.value?.characteristics.wits,
    multitasking: data.value?.characteristics.multitasking,
    management: data.value?.characteristics.management,
    planning: data.value?.characteristics.planning
  }
})

</script>

<style scoped>

</style>