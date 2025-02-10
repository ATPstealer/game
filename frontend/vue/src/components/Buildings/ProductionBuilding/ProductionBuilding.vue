<template>
  <div v-if="building?.buildingType?.title" class="flex flex-col gap-4">
    <MessageBlock
      v-if="messageData?.code"
      v-bind="messageData"
      class="md:w-1/3 self-center"
    />
    <div class="flex gap-4 justify-center w-full flex-col-reverse md:flex-row">
      <Card class="w-full max-w-md">
        <template #content>
          <div class="flex flex-col gap-4">
            <div class="font-bold">
              {{ t('buildings.production.startWork') }}:
            </div>
            <Dropdown
              v-model="duration"
              option-label="label"
              option-value="value"
              :options="timeValues"
            />
            <Button
              class-name="w-full"
              :disabled="!selectedBlueprint"
              :label="t('buildings.production.startWork')"
              @click="start"
            />
          </div>
        </template>
      </Card>
      <Card class="w-full max-w-md building-card">
        <template #content>
          <span class="font-bold">
            {{ t('buildings.one') }}:
          </span>
          <p><span>{{ t('common.type') }}</span>: {{ getTranslation({parent: 'buildings.types', child: building?.buildingType?.title}) }}</p>
          <p><span>{{ t('common.status') }}</span>: {{ building?.status }}</p>
          <p><span>{{ t('common.coordinates') }}</span>: {{ building?.x }}:{{ building?.y }}</p>
          <p><span>{{ t('common.level') }}</span> x <span>{{ t('common.square') }}</span>: {{ building?.level }}x{{ building?.square }}</p>
          <p><span>Занято места</span> {{ Math.round(building.squareInUse / (building.square * building.level) * 100) }}%</p>
        </template>
      </Card>
    </div>
    <h3 class="text-center">
      {{ t('buildings.production.chooseResource') }}
    </h3>
    <div v-if="blueprints?.length && resourceTypes?.length" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <ResourceCard
        v-for="blueprint in blueprints"
        :key="blueprint.id"
        :blueprint="blueprint"
        :building="building"
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
import { computed, ref, unref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { timeValues } from '@/components/Buildings/ProductionBuilding/constants'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import ResourceCard from '@/components/Resources/ResourceCard.vue'
import {
  type BuildingWithData,
  type JsonResult, type TimeDuration,
  useGetBuildingBlueprints,
  useGetResourceTypes,
  usePostBuildingStartWork
} from '@/gen'
import { getTranslation } from '@/utils/getTranslation'

// TODO: при изменении кол-ва работников/зп, при нажатии на крестик сейчас показываются значения из props.building,
// добавить перезапрос данных о building при изменении кол-ва работников/зп

interface Props {
  building: BuildingWithData | undefined;
}

const props = defineProps<Props>()

const duration = ref<TimeDuration>(3600 * 1000000000 as TimeDuration)
const selectedBlueprint = ref<number>(0)
const messageData = ref<JsonResult>()

const { t } = useI18n()

const { data: blueprintsQuery, suspense: awaitBlueprints  } = useGetBuildingBlueprints()
await awaitBlueprints()
const blueprints = computed(() => unref(blueprintsQuery)?.data?.filter(item => item.producedInId === props.building!.typeId) || [])

watch(blueprints, () => {
  if (blueprints.value.length === 1) {
    selectedBlueprint.value = blueprints.value[0].id
  }
})

const { data: resourceTypesQuery, suspense: awaitResourceTypes } = useGetResourceTypes()
await awaitResourceTypes()
const resourceTypes = computed(() => {
  return unref(resourceTypesQuery)?.data || []
})

const mutateStartProduction = usePostBuildingStartWork({
  mutation: {
    onSuccess: data => {
      messageData.value = data
    }
  }
})

const selectBlueprint = (event: number) => {
  if (selectedBlueprint.value === event) {
    selectedBlueprint.value = 0

    return
  }
  selectedBlueprint.value = event
}

const start = () => {
  messageData.value = {} as JsonResult

  const payload = {
    buildingId: props.building!._id,
    blueprintId: selectedBlueprint.value,
    duration: duration.value
  }

  mutateStartProduction.mutate({ data: { ...payload } })
}

</script>

<style scoped>
.building-card span {
  @apply capitalize;
}
</style>