<template>
  <div v-if="resourceTypes?.length" class="flex gap-2 flex-wrap">
    <span
      v-for="resource in resourceTypes"
      :key="resource._id"
      class="item item-hover"
      @click="chosen = resource"
      :class="{'bg-amber-100' : resource.id === chosen?.id }"
    >
      {{ resource.name }}
    </span>
  </div>

  <div class="flex justify-evenly items-center mt-10 h-[600px]">
    <div>
      <div v-if="parts.length" class="flex flex-col gap-2">
        <div
          v-for="partResource in parts"
          :key="partResource.id"
          class="flex relative"
        >
          <span
            class="item item-hover"
            @click="chosen = partResource"
          >
            {{ partResource.name }}
          </span>
          <span class="font-bold absolute -right-8 top-50">{{ getAmount(partResource) }}</span>
        </div>
      </div>
    </div>
    <Divider layout="vertical" />
    <div class="min-w-[120px]">
      <div v-if="chosen">
        <span>{{ chosen.name }}</span>
        <div v-if="bps.length" class="flex flex-col gap-2">
          <Divider />
          <span>Blueprints</span>
          <div
            v-for="bp in bps"
            :key="bp.id"
            class="item"
            @mouseover="bpHovered = bp"
            @mouseleave="bpHovered = {}"
          >
            {{ bp.name }}
          </div>
        </div>
        <div v-if="buildings.length" class="flex flex-col gap-2">
          <Divider />
          <span>Produced in</span>
          <div
            v-for="building in buildings"
            :key="building.id"
            class="item"
          >
            {{ building.title }}
          </div>
        </div>
      </div>
    </div>
    <Divider layout="vertical" />
    <divc class="min-w-[120px]">
      <div v-if="produce.length" class="flex flex-col gap-2 m-auto">
        <span
          class="item item-hover"
          v-for="producedResource in produce"
          :key="producedResource.id"
          @click="chosen = producedResource"
        >
          {{ producedResource.name }}
        </span>
      </div>
    </divc>
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

const chosen = ref<Resource>()
const bpHovered = ref<Blueprint>({} as Blueprint)

const { data: resourceTypes } = useGetData<Resource[]>('/resource/types')
const { data: blueprints } = useGetData<Blueprint[]>('/building/blueprints')
const { data: buildingsTypes } = useGetData<Building[]>('/building/types')

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

const buildings = computed(() => {
  if (!chosen.value) {
    return []
  }

  const bpUsed = blueprints.value.filter(item => {
    return item.producedResources.some(resource => {
      return resource.resourceId === chosen.value?.id
    })
  })

  const ids = bpUsed.map(item => item.producedInId)

  return buildingsTypes.value.filter(item => ids.includes(item.id))
})

const bps = computed(() => {
  if (!chosen.value) {
    return []
  }

  return blueprints.value.filter(item => {
    return item.producedResources.some(resource => {
      return resource.resourceId === chosen.value?.id
    })
  })
})

const getAmount = (resource: Resource) => {
  if (isEmpty(bpHovered.value)) {
    return
  }

  return bpHovered.value.usedResources.find(item => item.resourceId === resource.id).amount
}

</script>

<style scoped>
.item {
  @apply p-2 border border-solid cursor-pointer min-w-[120px] flex justify-center;
}

.item-hover {
  @apply hover:bg-amber-100;
}
</style>