<template>
  <div v-if="resourceTypes?.length" class="flex gap-2 flex-wrap">
    <span
      v-for="resource in resourceTypes"
      :key="resource.id"
      class="item item-hover"
      @click="chosen = resource"
      :class="{'bg-amber-100' : resource.id === chosen?.id }"
    >
      {{ resource.name }}
    </span>
  </div>

  <div class="flex justify-evenly items-center mt-10 h-[600px]">
    <div class="column-block">
      <span class="column-header">Produced from</span>
      <div v-if="computedData?.usedResources?.length" class="flex flex-col gap-2">
        <div
          v-for="partResource in computedData.usedResources"
          :key="partResource.id"
          class="flex relative"
        >
          <span
            class="item item-hover"
            @click="chosen = partResource"
          >
            {{ partResource.name }}
          </span>
          <span class="font-bold absolute -right-4 translate-x-full top-1/2 -translate-y-1/2">{{ getAmount(partResource) }}</span>
        </div>
      </div>
    </div>
    <Divider layout="vertical" />
    <div class="column-block">
      <span class="column-header">Resource</span>
      <div v-if="chosen">
        <span>{{ chosen.name }}</span>
        <div v-if="computedData?.blueprints?.length" class="flex flex-col gap-2">
          <Divider />
          <span>Blueprints</span>
          <div
            v-for="bp in computedData.blueprints"
            :key="bp.id"
            class="item item-hover"
            @mouseover="bpHovered = bp"
            @mouseleave="bpHovered = {} as Blueprint"
          >
            {{ bp.name }}
          </div>
        </div>
        <div v-if="computedData?.buildings?.length" class="flex flex-col gap-2">
          <Divider />
          <span>Produced in</span>
          <div
            v-for="building in computedData.buildings"
            :key="building.id"
            class="item !cursor-default"
          >
            {{ building.title }}
          </div>
        </div>
      </div>
    </div>
    <Divider layout="vertical" />
    <div class="column-block">
      <span class="column-header">Used in</span>
      <div v-if="computedData?.producedResources?.length" class="flex flex-col gap-2 m-auto">
        <span
          class="item item-hover"
          v-for="producedResource in computedData.producedResources"
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
import isEmpty from 'lodash/isEmpty'
import Divider from 'primevue/divider'
import { computed, ref } from 'vue'
import { useGetData } from '@/composables/useGetData'
import type { Blueprint, Building } from '@/types/Buildings/index.interface'
import type { Resource } from '@/types/Resources/index.interface'

interface Pipeline {
  usedResources: Resource[];
  producedResources: Resource[];
  blueprints: Blueprint[];
  buildings: Building[];
}

const chosen = ref<Resource>()
const bpHovered = ref<Blueprint>({} as Blueprint)

const { data: resourceTypes } = useGetData<Resource[]>('/resource/types')
const { data: blueprints } = useGetData<Blueprint[]>('/building/blueprints')
const { data: buildingsTypes } = useGetData<Building[]>('/building/types')

const computedData = computed<Pipeline>(() => {
  if (!chosen.value) {
    return {} as Pipeline
  }

  const bpUsed = blueprints.value.filter(item => {
    return item.usedResources.some(resource => {
      return resource.resourceId === chosen.value?.id
    })
  })

  const bpProd = blueprints.value.filter(item => {
    return item.producedResources.some(resource => {
      return resource.resourceId === chosen.value?.id
    })
  })

  const used = uniq(bpUsed.map(item => item.producedResources.map(i => i.resourceId)).flat())
  const prod = uniq(bpProd.map(item => item.usedResources.map(i => i.resourceId)).flat())
  const buildings = bpProd.map(item => item.producedInId)

  return {
    usedResources: resourceTypes.value.filter(item => prod.includes(item.id)),
    producedResources: resourceTypes.value.filter(item => used.includes(item.id)),
    buildings: buildingsTypes.value.filter(item => buildings.includes(item.id)),
    blueprints: bpProd
  }
})

const getAmount = (resource: Resource) => {
  if (isEmpty(bpHovered.value)) {
    return
  }

  return bpHovered.value.usedResources.find(item => item.resourceId === resource.id)?.amount
}

</script>

<style scoped>
.item {
  @apply p-2 border border-solid cursor-pointer min-w-[120px] flex justify-center;
}

.item-hover {
  @apply hover:bg-amber-100;
}

.column-block {
  @apply min-w-[200px] relative h-full flex items-center justify-center;
}

.column-header {
  @apply absolute top-0 text-lg font-bold;
}
</style>