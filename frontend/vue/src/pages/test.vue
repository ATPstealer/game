<template>
  <div v-if="resourceTypes?.length" class="flex gap-2 flex-wrap">
    <span
      v-for="resource in resourceTypes"
      :key="resource._id"
      class="item"
      @click="chosen = resource"
    >
      {{ resource.name }}
    </span>
  </div>

  <div class="flex justify-between mt-10 h-[500px]">
    <div class="self-center">
      <div v-if="parts.length" class="flex flex-col gap-2">
        <span
          class="item"
          v-for="partResource in parts"
          :key="partResource.id"
          @click="chosen = partResource"
        >
          {{ partResource.name }}
        </span>
      </div>
    </div>
    <Divider layout="vertical" />
    <div class="self-center">
      <div v-if="chosen">
        {{ chosen.name }}
      </div>
    </div>
    <Divider layout="vertical" />
    <div class="self-center">
      <div v-if="produce.length" class="flex flex-col gap-2">
        <span
          class="item"
          v-for="producedResource in produce"
          :key="producedResource.id"
          @click="chosen = producedResource"
        >
          {{ producedResource.name }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { uniq } from 'lodash'
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

  const bpUsed = blueprints.value.filter(item => {
    return item.usedResources.some(resource => {
      return resource.resourceId === chosen.value?.id
    })
  })

  const used = bpUsed.map(item => item.producedResources.map(i => i.resourceId)).flat()

  const ids = uniq(used)

  return resourceTypes?.value.filter(item => ids.includes(item.id))
})

const parts  = computed(() => {
  if (!chosen.value) {
    return []
  }

  const bpUsed = blueprints.value.filter(item => {
    return item.producedResources.some(resource => {
      return resource.resourceId === chosen.value?.id
    })
  })

  const used = bpUsed.map(item => item.usedResources.map(i => i.resourceId)).flat()

  const ids = uniq(used)

  return resourceTypes?.value.filter(item => ids.includes(item.id))
})

</script>

<style scoped>
.item {
  @apply p-2 border border-solid cursor-pointer hover:bg-amber-100;
}
</style>