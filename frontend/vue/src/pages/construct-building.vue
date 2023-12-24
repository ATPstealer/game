<template>
  <div class="flex flex-col items-center justify-center mt-5 space-y-10">
    <span class="text-5xl font-bold">Building Builder</span>
    <MessageBlock v-if="message" :message="message" />
    <Card>
      <template #content>
        <div class="flex flex-col gap-4">
          <p class="label">
            Where:
          </p>
          <div class="flex gap-4 w-full">
            <div class="flex flex-col flex-1">
              <label for="x" class="label">X:</label>
              <InputNumber
                show-buttons
                v-model="x"
                input-id="x"
              />
            </div>
            <div class="flex flex-col flex-1">
              <label for="y" class="label">Y:</label>
              <InputNumber
                show-buttons
                v-model="y"
                input-id="y"
              />
            </div>
          </div>
          <div class="flex gap-4 w-full">
            <div class="flex-1">
              <label for="type_id" class="label">Build:</label>
              <Dropdown
                :options="buildingTypes"
                v-model="buildingType"
                option-label="title"
                class="w-full"
              />
            </div>
            <div class="flex-1 self-center text-xl">
              {{ buildingType.description }}
            </div>
          </div>
          <div class="flex flex-col">
            <label for="square" class="label">Square:</label>
            <InputNumber
              show-buttons
              v-model="square"
              input-id="square"
            />
          </div>
          <div class="flex space-x-3">
            <div class="flex flex-col">
              <div class="flex gap-4">
                <p class="label">
                  Cost:
                </p>
                <p class="text-xl underline">
                  {{ buildingType.cost * square }}$
                </p>
              </div>
              <div class="flex gap-4">
                <p class="label">
                  Time:
                </p>
                <p class="text-xl underline">
                  {{ formatDuration(buildingType.buildTime * square / 1000000000) }}
                </p>
              </div>
            </div>
          </div>

          <Button
            label="Construct"
            class="bg-indigo-500"
            @click="construct"
          />
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Card from 'primevue/card'
import Dropdown from 'primevue/dropdown'
import InputNumber from 'primevue/inputnumber'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useBuildings } from '@/composables/useBuildings'
import { useGetData } from '@/composables/useGetData'
import type { Message } from '@/types'
import type { BuildingType } from '@/types/Buildings/index.interface'
import { formatDuration } from '@/utils/formatDuration'

const { query } = useRoute()

const x = ref<number>(Number(query.x))
const y = ref<number>(Number(query.y))
const buildingType = ref<BuildingType>({} as BuildingType)
const square = ref<number>(10)
const message = ref<Message | null>(null)

const { data: buildingTypes, onFetchResponse } = useGetData<BuildingType[]>('/building/types')
onFetchResponse(() => {
  buildingType.value = buildingTypes.value[0]
})

const { constructBuilding } = useBuildings()

const construct = () => {
  message.value = null

  const payload = {
    x: x.value,
    y: y.value,
    typeId: buildingType.value.id,
    square: square.value
  }

  const { dataMessage, onFetchFinally } = constructBuilding(payload)
  onFetchFinally(() => {
    message.value = dataMessage.value
  })
}

</script>

<style scoped>
.label {
  @apply block text-xl font-bold mb-2;
}
</style>