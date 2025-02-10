<template>
  <PipelinesTemplate :headers="headers">
    <template #items>
      <div v-if="buildingTypes?.length" class="flex gap-2 flex-wrap">
        <span
          v-for="building in buildingTypes.filter(item => item.buildingGroup === 'Production')"
          :key="building.id"
          class="pipeline-item pipeline-item-hover"
          :class="{'bg-amber-100' : building.id === chosen?.id }"
          @click="chosen = building"
        >
          {{ t(`buildings.types.${building?.title.toLowerCase()}`) }}
        </span>
      </div>
    </template>
    <template #first-column>
      <div v-if="computedData?.usedBuildings?.length" class="flex flex-col gap-2">
        <div v-for="item in computedData.usedBuildings" :key="item.id">
          <span class="pipeline-item pipeline-item-hover" @click="chosen = item">{{ t(`buildings.types.${item.title.toLowerCase()}`) }}</span>
        </div>
      </div>
    </template>
    <template #second-column>
      <div v-if="chosen" class="flex flex-col gap-2">
        <span class="text-xl underline text-center">{{ t(`buildings.types.${chosen?.title.toLowerCase()}`) }}</span>
      </div>
    </template>
    <template #third-column>
      <div v-if="computedData?.prodBuildings?.length" class="flex flex-col gap-2 items">
        <div v-for="item in computedData.prodBuildings" :key="item.id">
          <span class="pipeline-item pipeline-item-hover" @click="chosen = item">{{ t(`buildings.types.${item.title.toLowerCase()}`) }}</span>
        </div>
      </div>
    </template>
  </PipelinesTemplate>
</template>

<script setup lang="ts">
import uniq from 'lodash/uniq'
import { computed, ref, toRef } from 'vue'
import { useI18n } from 'vue-i18n'
import PipelinesTemplate from '@/components/Pipelines/PipelinesTemplate.vue'
import type { BuildingType, Blueprint, ResourceType } from '@/gen'

interface Props {
  blueprints: Blueprint[];
  buildingTypes: BuildingType[];
  resourceTypes: ResourceType[];
}

interface Pipeline {
  usedBuildings: BuildingType[];
  prodBuildings: BuildingType[];
}

const props = defineProps<Props>()
const resourceTypes = toRef(props, 'resourceTypes')
const blueprints = toRef(props, 'blueprints')
const buildingTypes = toRef(props, 'buildingTypes')

const { t } = useI18n()

const chosen = ref<BuildingType>(buildingTypes?.value?.filter(item => item.buildingGroup === 'Production')[0] || {} as BuildingType)
// TODO: добавить переводы
const headers = ['Потребляет ресурсы из', t('common.building'), 'Производит ресурсы для']

const computedData = computed<Pipeline>(() => {
  if (!chosen.value) {
    return {} as Pipeline
  }

  const buildingBps = blueprints.value.filter(item => item.producedInId === chosen.value?.id)

  const usedResources = uniq(buildingBps.flatMap(bp => bp.usedResources.map(item => item.resourceId)))
  const bpsUsed = blueprints.value.filter(item => item.producedResources.some(resource => usedResources.includes(resource.resourceId)))
  const buildingsUsed = uniq(bpsUsed.flatMap(bp => bp.producedInId))

  const producedResources = uniq(buildingBps.flatMap(bp => bp.producedResources.map(item => item.resourceId)))
  const bpProd = blueprints.value.filter(item => item.usedResources.some(resource => producedResources.includes(resource.resourceId)))
  const buildingsProd = uniq(bpProd.flatMap(bp => bp.producedInId))

  return {
    usedBuildings: buildingTypes.value.filter(item => buildingsUsed.includes(item.id) && item.id !== chosen.value?.id).filter(item => item.buildingGroup === 'Production'),
    prodBuildings: buildingTypes.value.filter(item => buildingsProd.includes(item.id) && item.id !== chosen.value?.id).filter(item => item.buildingGroup === 'Production')
  }
})
</script>

<style scoped>

</style>