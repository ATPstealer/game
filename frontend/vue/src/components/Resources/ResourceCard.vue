<template>
  <Card
    class="w-full relative"
    :class="{'bg-blue-200': selectedBlueprint === blueprint.id}"
    @click="emits('select', blueprint.id)"
  >
    <template #title>
      <span class="text-xl">{{ blueprint.name }}</span>
    </template>
    <template #content>
      <div class="flex flex-col gap-2">
        <p class="font-bold">
          {{ t('buildings.production.produce') }}:
        </p>
        <div>
          <p
            v-for="resource in blueprint.producedResources"
            :key="resource.resourceId"
            class="ml-4"
          >
            {{ t(`resources.types.${findResourceName(resourceTypes, +resource.resourceId)?.toLowerCase()}`) }} {{ resource.amount }}
          </p>
        </div>
        <p class="font-bold">
          {{ t('buildings.production.use') }}:
        </p>
        <div>
          <p
            v-for="resource in blueprint.usedResources"
            :key="resource.resourceId"
            class="ml-4"
          >
            {{ t(`resources.types.${findResourceName(resourceTypes, +resource.resourceId)?.toLowerCase()}`) }} {{ resource.amount }}
          </p>
        </div>
        <p><span class="font-bold">{{ t('buildings.production.cycle') }}</span>: {{ getCycling() }}s</p>
        <img
          :alt="buildingIcon"
          class="max-h-[64px] absolute top-4 right-4"
          :src="getMinioURL(`/resource/${buildingIcon}`)"
        />
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import { useI18n } from 'vue-i18n'
import type { Blueprint, Building, EquipmentEffect } from '@/types/Buildings/index.interface'
import type { ResourceType } from '@/types/Resources/index.interface'
import { findResourceName } from '@/utils/findResourceName'
import { getMinioURL } from '@/utils/getMinioURL'

interface Props {
  blueprint: Blueprint;
  resourceTypes: ResourceType[];
  selectedBlueprint: number | undefined;
  building: Building;
}

const { t } = useI18n()
const props = defineProps<Props>()

const emits = defineEmits<{(e: 'select', value: number)}>()

const buildingIcon = findResourceName(props.resourceTypes, +props.blueprint.producedResources[0].resourceId )

const getCycling = () => {
  console.log(props.building.equipmentEffect)
  if (!props.building.equipmentEffect) {
    return Math.round(props.blueprint.productionTime * props.building.hiringNeeds / props.building.workers / 1000000000)
  }
  const effectValue = props.building.equipmentEffect.reduce((acc: { bp: number; all: number }, effect: EquipmentEffect) => {
    if (effect.blueprintId === props.blueprint.id) {
      acc.bp += effect.value
    }
    if (effect.blueprintId === 0) {
      acc.all += effect.value
    }

    return acc
  }, { bp: 0, all: 0 })

  const effectiveness = props.building.workers + effectValue.bp + effectValue.all

  if (effectiveness === 0) {
    return '∞'
  }

  return Math.round(props.blueprint.productionTime * props.building.hiringNeeds / (effectiveness ) / 1000000000)
}
</script>

<style scoped>

</style>
