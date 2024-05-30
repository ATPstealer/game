<template>
  <PipelinesTemplate :headers="headers">
    <template #items>
      <div v-if="resourceTypes?.length" class="flex gap-2 flex-wrap">
        <span
          v-for="resource in resourceTypes"
          :key="resource.id"
          class="pipeline-item pipeline-item-hover"
          :class="{'bg-amber-100' : resource.id === chosen?.id }"
          @click="chosen = resource"
        >
          {{ t(`resources.types.${resource.name.toLowerCase()}`) }}
        </span>
      </div>
    </template>
    <template #first-column>
      <div v-if="computedData?.usedResources?.length" class="flex flex-col gap-2">
        <div
          v-for="partResource in computedData.usedResources"
          :key="partResource.id"
          class="flex relative"
        >
          <span
            class="pipeline-item pipeline-item-hover"
            @click="chosen = partResource"
          >
            {{ t(`resources.types.${partResource.name.toLowerCase()}`) }}
          </span>
          <span class="font-bold absolute -right-4 translate-x-full top-1/2 -translate-y-1/2">{{ getAmount(partResource) }}</span>
        </div>
      </div>
    </template>
    <template #second-column>
      <div v-if="chosen" class="flex flex-col gap-2">
        <span class="text-xl underline text-center">{{ t(`resources.types.${chosen.name.toLowerCase()}`) }}</span>
        <div v-if="computedData?.blueprints?.length" class="flex flex-col gap-2">
          <Divider />
          <span>Blueprints</span>
          <div
            v-for="bp in computedData.blueprints"
            :key="bp.id"
            class="pipeline-item pipeline-item-hover"
            @mouseleave="bpHovered = {} as Blueprint"
            @mouseover="bpHovered = bp"
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
            class="pipeline-item !cursor-default"
          >
            {{ t(`buildings.types.${building.title.toLowerCase()}`) }}
          </div>
        </div>
      </div>
    </template>
    <template #third-column>
      <div
        v-if="computedData?.producedResources?.length"
        class="flex flex-col gap-2 m-auto"
        :class="{'!grid !grid-cols-2': computedData?.producedResources?.length > 10}"
      >
        <span
          v-for="producedResource in computedData.producedResources"
          :key="producedResource.id"
          class="pipeline-item pipeline-item-hover"
          @click="chosen = producedResource"
        >
          {{ t(`resources.types.${producedResource.name.toLowerCase()}`) }}
        </span>
      </div>
    </template>
  </PipelinesTemplate>
</template>

<script setup lang="ts">
import { uniq } from 'lodash'
import isEmpty from 'lodash/isEmpty'
import Divider from 'primevue/divider'
import { computed, ref, toRef } from 'vue'
import { useI18n } from 'vue-i18n'
import PipelinesTemplate from '@/components/Pipelines/PipelinesTemplate.vue'
import type { Blueprint, Building } from '@/types/Buildings/index.interface'
import type { Resource } from '@/types/Resources/index.interface'

interface Props {
  resourceTypes: Resource[];
  blueprints: Blueprint[];
  buildingsTypes: Building[];
}
interface Pipeline {
  usedResources: Resource[];
  producedResources: Resource[];
  blueprints: Blueprint[];
  buildings: Building[];
}

const props = defineProps<Props>()
const resourceTypes = toRef(props, 'resourceTypes')
const blueprints = toRef(props, 'blueprints')
const buildingsTypes = toRef(props, 'buildingsTypes')

const chosen = ref<Resource>(resourceTypes?.value?.[0] || {} as Resource)
const bpHovered = ref<Blueprint>({} as Blueprint)

const { t } = useI18n()
// TODO: добавить переводы
const headers = ['Produced from', 'Resource', 'Used in']

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

</style>