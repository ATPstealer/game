<template>
  <Layout :show-options-prop="false">
    <div class="flex flex-col gap-4 items-center justify-center">
      <h2>{{ t('buildings.builder') }}</h2>
      <MessageBlock
        v-if="messageData?.code"
        v-bind="messageData"
      />
      <Card v-if="isSuccess">
        <template #content>
          <div class="flex flex-col gap-4">
            <p class="label !mb-0">
              {{ t('common.where') }}:
            </p>
            <div class="flex gap-4 w-full">
              <div class="coordinate">
                <label class="label" for="x">X:</label>
                <InputNumber
                  v-model="x"
                  input-class="max-w-[70px] md:max-w-[unset]"
                  input-id="x"
                  show-buttons
                />
              </div>
              <div class="coordinate">
                <label class="label" for="y">Y:</label>
                <InputNumber
                  v-model="y"
                  input-class="max-w-[70px] md:max-w-[unset]"
                  input-id="y"
                  show-buttons
                />
              </div>
            </div>
            <div class="flex flex-col md:flex-row gap-4 w-full">
              <div class="flex-1">
                <label class="label">{{ t('common.build') }}:</label>
                <Dropdown
                  v-model="buildingType"
                  class="w-full"
                  option-label="title"
                  :options="buildingTypes"
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
              <label class="label" for="square">{{ t('common.square') }}:</label>
              <InputNumber
                v-model="square"
                input-id="square"
                show-buttons
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
              class="bg-indigo-500"
              :label="t('common.construct')"
              @click="construct"
            />
          </div>
        </template>
      </Card>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { useMutation, useQuery } from '@tanstack/vue-query'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Dropdown from 'primevue/dropdown'
import InputNumber from 'primevue/inputnumber'
import { type Ref, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import {
  type models_BuildingType,
  type models_ConstructBuildingPayload,
  type PostApiV2BuildingConstructResponse
} from '@/api'
import { getApiV2BuildingTypesOptions, postApiV2BuildingConstructMutation } from '@/api/@tanstack/vue-query.gen'
import Layout from '@/components/Common/Layout.vue'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import type { BackData } from '@/types'
import type { BuildingType } from '@/types/Buildings/index.interface'
import { formatDuration } from '@/utils/formatDuration'
import { getTranslation } from '@/utils/getTranslation'

const { query } = useRoute()

const x = ref<number>(Number(query.x))
const y = ref<number>(Number(query.y))
const buildingType = ref<models_BuildingType>({} as models_BuildingType)
const square = ref<number>(10)
const messageData = ref<PostApiV2BuildingConstructResponse>()

const { t } = useI18n()

const { data: buildingTypes, isSuccess } = useQuery({
  ...getApiV2BuildingTypesOptions()
})

const constructBuilding = useMutation({
  ...postApiV2BuildingConstructMutation(),
  onSuccess: (data) => {
    messageData.value = data
  }
})

watch(isSuccess, () => {
  if (isSuccess.value) {
    buildingType.value = buildingTypes.value[0]
  }
})

const construct = () => {
  messageData.value = {} as BackData

  const payload: models_ConstructBuildingPayload = {
    x: x.value,
    y: y.value,
    typeId: buildingType.value.id,
    square: square.value
  }

  constructBuilding.mutate({ body: { ...payload } })
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
