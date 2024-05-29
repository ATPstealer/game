<template>
  <div v-if="resourceTypes?.length" class="flex gap-2 flex-wrap">
    <span
      v-for="resource in resourceTypes"
      :key="resource._id"
      class="p-2 border border-solid cursor-pointer hover:bg-amber-100"
      @click="chosen = resource"
    >
      {{ resource.name }}
    </span>
  </div>

  <div class="flex justify-between mt-10 h-[500px]">
    <div>1</div>
    <Divider layout="vertical" />
    <div class="self-center">
      <div v-if="chosen">
        {{ chosen.name }}
      </div>
    </div>
    <Divider layout="vertical" />
    <div>3</div>
  </div>

  {{ produce }}
</template>

<script setup lang="ts">
import Divider from 'primevue/divider'
import { computed, ref } from 'vue'
import { useGetData } from '@/composables/useGetData'
import type { Blueprint } from '@/types/Buildings/index.interface'
import type { Resource } from '@/types/Resources/index.interface'

const chosen = ref<Resource>()

const { data: resourceTypes } = useGetData<Resource[]>('/resource/types')
const { data: blueprints } = useGetData<Blueprint[]>('/building/blueprints')

const produce = computed(() => {
  if (!chosen.value) {
    return []
  }

  const bp = blueprints.value.filter(item => {
    return item.usedResources.some(resource => {
      return resource.resourceId === chosen.value?.id
    })
  })

  const result = bp.map(item => {

  })

  return []
})

const parts  = computed(() => {
  if (!chosen.value) {
    return []
  }

  const bp = blueprints.value.filter(item => {
    console.log(item)

    return item.producedResources.some(resource => resource.resourceId === chosen.value?._id)
  })
  console.log(bp)

  return []
})

</script>

<style scoped>

</style>