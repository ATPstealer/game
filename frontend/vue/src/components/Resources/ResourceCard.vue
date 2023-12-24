<template>
  <Card
    @click="emits('select', blueprint.id)"
    :class="{'bg-blue-200': selectedBlueprint === blueprint.id}"
    class="w-full relative"
  >
    <template #title>
      <span class="text-xl">{{ blueprint.name }}</span>
    </template>
    <template #content>
      <div class="flex flex-col gap-2">
        <p class="font-bold">
          Produce:
        </p>
        <div>
          <p
            v-for="resource in blueprint.producedResources"
            :key="resource.resourceId"
            class="ml-4"
          >
            {{ findResourceName(resourceTypes, resource.resourceId) }} {{ resource.amount }}
          </p>
        </div>
        <p class="font-bold">
          Use:
        </p>
        <div>
          <p
            v-for="resource in blueprint.usedResources"
            :key="resource.resourceId"
            class="ml-4"
          >
            {{ findResourceName(resourceTypes, resource.resourceId) }} {{ resource.amount }}
          </p>
        </div>
        <p><span class="font-bold">1 cycle time</span>: {{ blueprint.productionTime / 1000000000 }}s</p>
        <img
          class="max-h-[64px] absolute top-4 right-4"
          :src="getMinioURL(`/resource/${buildingIcon}`)"
          :alt="buildingIcon"
        />
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import type { Blueprint } from '@/types/Buildings/index.interface'
import type { ResourceType } from '@/types/Resources/index.interface'
import { findResourceName } from '@/utils/findResourceName'
import { getMinioURL } from '@/utils/getMinioURL'

interface Props {
  blueprint: Blueprint;
  resourceTypes: ResourceType[];
  selectedBlueprint: number | undefined;
}

const props = defineProps<Props>()

const emits = defineEmits<{(e: 'select', value: number)}>()

const buildingIcon = findResourceName(props.resourceTypes, props.blueprint.producedResources[0].resourceId )
</script>

<style scoped>

</style>