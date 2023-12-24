<template>
  <div class="flex flex-col gap-4">
    <h1 class="text-center">
      Production
    </h1>
    <MessageBlock
      :message="message"
      v-if="message"
      class="md:w-1/3 self-center"
    />
    <div class="flex gap-4 justify-center w-full flex-col-reverse md:flex-row">
      <Card class="w-full max-w-md">
        <template #content>
          <div class="flex flex-col gap-4">
            <div class="font-bold">
              Start work:
            </div>
            <Dropdown
              v-model="duration"
              :options="timeValues"
              option-label="label"
              option-value="value"
            />
            <Button
              :disabled="!selectedBlueprint"
              @click="start"
              label="Start work"
              class-name="w-full"
            />
          </div>
        </template>
      </Card>
      <Card class="w-full max-w-md">
        <template #content>
          <span class="font-bold">
            Building:
          </span>
          <p>Type: {{ building?.title }}</p>
          <p>Status: {{ building?.status }}</p>
          <p>Coordinates: {{ building?.x }}:{{ building?.y }}</p>
          <p>Level x Square: {{ building?.level }}x{{ building?.square }}</p>
        </template>
      </Card>
    </div>
    <h3 class="text-center">
      Choose resource for creating:
    </h3>
    <div v-if="blueprints?.length && resourceTypes?.length" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <ResourceCard
        v-for="blueprint in blueprints"
        :key="blueprint.id"
        :blueprint="blueprint"
        :resource-types="resourceTypes"
        :selected-blueprint="selectedBlueprint"
        @select="selectBlueprint"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Card from 'primevue/card'
import Dropdown from 'primevue/dropdown'
import { ref } from 'vue'
import { timeValues } from '@/components/Buildings/ProductionBuilding/constants'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import ResourceCard from '@/components/Resources/ResourceCard.vue'
import { useBuildings } from '@/composables/useBuildings'
import { useGetData } from '@/composables/useGetData'
import type { Message } from '@/types'
import type { Blueprint, Building } from '@/types/Buildings/index.interface'
import type { ResourceType } from '@/types/Resources/index.interface'

interface Props {
  building: Building;
}

const props = defineProps<Props>()

const blueprints = ref<Blueprint[]>([])
const duration = ref<number>(3600)
const selectedBlueprint = ref<number>(0)
const message = ref<Message | null>(null)

const { data, onFetchResponse } = useGetData('/building/blueprints')
onFetchResponse(() => {
  blueprints.value = data.value.filter(item => item.producedInId === props.building.typeId)
  if (blueprints.value.length === 1) {
    selectedBlueprint.value = blueprints.value[0].id
  }
})

const { data: resourceTypes } = useGetData<ResourceType[]>('/resource/types')

const selectBlueprint = (event: number) => {
  if (selectedBlueprint.value === event) {
    selectedBlueprint.value = 0

    return
  }
  selectedBlueprint.value = event
}

const start = () => {
  const { startProduction } = useBuildings()

  const payload = {
    buildingId: props.building.id,
    blueprintId: selectedBlueprint.value,
    duration: duration.value
  }

  const { dataMessage, onFetchResponse } = startProduction(payload)
  onFetchResponse(() => {
    message.value = dataMessage.value
  })
}

</script>

<style scoped>

</style>