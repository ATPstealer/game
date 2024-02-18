<template>
  <MessageBlock
    v-if="message"
    :message="message"
    class="mb-4"
  />
  <div class="m-0">
    <p class="font-bold">
      {{ t(`buildings.hiring.name`) }}
    </p>
    <p class="italic" v-if="building.onStrike">
      {{ t(`buildings.hiring.on strike`) }}
    </p>
    <p>
      {{ t(`buildings.hiring.max workers`) }}: {{ building.maxWorkers*building.level*building.square }}
    </p>
    <p>
      {{ t(`buildings.hiring.workers`) }} / {{ t(`buildings.hiring.hiring needs`) }}:
    </p>
    <div class="flex items-center gap-1 h-10">
      <p class="whitespace-nowrap">
        {{ props.building.workers }} /
      </p>
      <span
        v-if="!editHiringNeeds"
        @click="editHiringNeeds = true; editSalary = false; message = null"
        class="font-bold text-blue-500 hover:text-blue-700 cursor-pointer"
      >
        {{ hiringNeeds ? hiringNeeds : t(`buildings.hiring.set`) }}
      </span>
      <div v-else class="flex items-center gap-1">
        <InputNumber
          v-model="hiringNeeds"
          show-buttons
          input-class="!p-2 !w-1/2 !min-w-[60px]"
          class="!w-1/2"
        />
        <Button
          icon="pi pi-check"
          size="small"
          severity="info"
          @click="setHiringData(hiringNeeds, 'needs')"
        />
        <Button
          size="small"
          icon="pi pi-times"
          severity="danger"
          @click="editHiringNeeds = false"
        />
      </div>
    </div>
    <div>
      <p>
        {{ t(`buildings.hiring.salary`) }}:
      </p>
      <span
        v-if="!editSalary"
        @click="editSalary = true; editHiringNeeds = false; message = null"
        class="font-bold text-blue-500 hover:text-blue-700 cursor-pointer"
      >
        {{ salary ? moneyFormat(salary) : t(`buildings.hiring.set`) }}
      </span>
      <div v-else class="flex items-center gap-1">
        <InputNumber
          v-model="salary"
          show-buttons
          :step="10"
          input-class="!p-2 !w-2/3"
          class="!w-2/3"
        />
        <Button
          icon="pi pi-check"
          size="small"
          severity="info"
          @click="setHiringData(salary, 'salary')"
        />
        <Button
          size="small"
          icon="pi pi-times"
          severity="danger"
          @click="editSalary = false"
        />
      </div>
    </div>
    <div>
      {{ t(`buildings.coefficient efficiency`) }}:
      {{ (building.workers / building.maxWorkers).toFixed(2) }}
    </div>
    <hr />
    <div v-if="!isMapFetching">
      {{ t(`buildings.hiring.average salary`) }}:
      {{ moneyFormat(getAverageSalary()) }}
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useBuildings } from '@/composables/useBuildings'
import { useMap } from '@/composables/useMap'
import { DataMessage } from '@/types'
import { Building } from '@/types/Buildings/index.interface'
import { moneyFormat } from '@/utils/moneyFormat'

interface Props {
  building: Building;
}

type HiringOptions = 'salary' | 'needs'
const props = defineProps<Props>()

const editSalary = ref<boolean>(false)
const editHiringNeeds = ref<boolean>(false)
const salary = ref<number>(props.building.salary)
const hiringNeeds = ref<number>(props.building.hiringNeeds)
const message = ref<DataMessage | null>(null)

const { setHiring } = useBuildings()
const { t } = useI18n()
const { getMap } = useMap()
const { data: map, isFetching: isMapFetching } = getMap()
const getAverageSalary = () => {
  if (map.value?.length) {
    return  map.value.filter(item => item.x === props.building.x && item.y === props.building.y)[0].averageSalary
  }

  return 0
}

const setHiringData = (value, option: HiringOptions) => {
  const isSalary = option === 'salary'

  const payload = {
    buildingId: props.building._id,
    salary: isSalary ? value : salary.value,
    hiringNeeds: !isSalary ? value : hiringNeeds.value
  }

  const { dataMessage: dataMessageSalary, onFetchResponse: onFetchResponseSalary } = setHiring(payload)

  onFetchResponseSalary(() => {
    message.value = dataMessageSalary.value

    editHiringNeeds.value = false
    editSalary.value = false
  })
}

</script>
