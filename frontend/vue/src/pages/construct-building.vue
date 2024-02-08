<template>
  <Layout :show-options-prop="false">
    <div class="flex flex-col gap-4 items-center justify-center">
      <h2>{{ t('buildings.builder') }}</h2>
      <MessageBlock v-if="message" :message="message" />
      <Card v-if="buildingTypes">
        <template #content>
          <div class="flex flex-col gap-4">
            <p class="label !mb-0">
              {{ t('common.where') }}:
            </p>
            <div class="flex gap-4 w-full">
              <div class="coordinate">
                <label for="x" class="label">X:</label>
                <InputNumber
                  show-buttons
                  v-model="x"
                  input-id="x"
                  input-class="max-w-[70px] md:max-w-[unset]"
                />
              </div>
              <div class="coordinate">
                <label for="y" class="label">Y:</label>
                <InputNumber
                  show-buttons
                  v-model="y"
                  input-id="y"
                  input-class="max-w-[70px] md:max-w-[unset]"
                />
              </div>
            </div>
            <div class="flex flex-col md:flex-row gap-4 w-full">
              <div class="flex-1">
                <label class="label">{{ t('common.build') }}:</label>
                <Dropdown
                  :options="buildingTypes"
                  v-model="buildingType"
                  option-label="title"
                  class="w-full"
                >
                  <template #option="{option}: {option: BuildingType}">
                    {{ getTranslation({parent: 'buildings.types', child: option.title}) }}
                  </template>
                  <template #value="{value}: {value: BuildingType}">
                    {{ getTranslation({parent: 'buildings.types', child: value.title}) }}
                  </template>
                </Dropdown>
              </div>
              <p class="flex-1 self-start md:self-end text-xl md:pb-3">
                {{ getTranslation({parent: 'buildings.typesDescriptions', child: buildingType.title}) }}
              </p>
            </div>
            <div class="flex flex-col">
              <label for="square" class="label">{{ t('common.square') }}:</label>
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
                    {{ t('common.cost') }}:
                  </p>
                  <p class="text-xl underline">
                    {{ buildingType.cost * square }}$
                  </p>
                </div>
                <div class="flex gap-4">
                  <p class="label">
                    {{ t('common.time') }}:
                  </p>
                  <p class="text-xl underline">
                    {{ formatDuration(buildingType.buildTime * square / 1000000000) }}
                  </p>
                </div>
              </div>
            </div>
            <Button
              :label="t('common.construct')"
              class="bg-indigo-500"
              @click="construct"
            />
          </div>
        </template>
      </Card>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Card from 'primevue/card'
import Dropdown from 'primevue/dropdown'
import InputNumber from 'primevue/inputnumber'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import Layout from '@/components/Common/Layout.vue'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useBuildings } from '@/composables/useBuildings'
import { useGetData } from '@/composables/useGetData'
import type { DataMessage } from '@/types'
import type { BuildingType } from '@/types/Buildings/index.interface'
import { formatDuration } from '@/utils/formatDuration'
import { getTranslation } from '@/utils/getTranslation'

const { query } = useRoute()

const x = ref<number>(Number(query.x))
const y = ref<number>(Number(query.y))
const buildingType = ref<BuildingType>({} as BuildingType)
const square = ref<number>(10)
const message = ref<DataMessage | null>(null)

const { constructBuilding } = useBuildings()
const { t } = useI18n()

const { data: buildingTypes, onFetchResponse } = useGetData<BuildingType[]>('/building/types')

onFetchResponse(() => {
  buildingType.value = buildingTypes.value[0]
})

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

.coordinate {
  @apply flex flex-row md:flex-col items-center md:items-start gap-4 md:gap-0 flex-1;
}
</style>
