<template>
  <MessageBlock
    v-if="messageData?.code"
    v-bind="messageData"
    class="mb-4"
  />
  <div v-if="building" class="m-0">
    <p class="font-bold">
      {{ t(`buildings.hiring.name`) }}
    </p>
    <p v-if="building.onStrike" class="italic">
      {{ t(`buildings.hiring.on strike`) }}
    </p>
    <p>
      {{ t(`buildings.hiring.max workers`) }}: {{ building?.buildingType?.workers * building?.level * building?.square }}
    </p>
    <p>
      {{ t(`buildings.hiring.workers`) }} / {{ t(`buildings.hiring.hiring needs`) }}:
    </p>
    <div class="flex items-center gap-1 h-10">
      <p class="whitespace-nowrap">
        {{ building.workers }} /
      </p>
      <span
        v-if="!editHiringNeeds"
        class="font-bold text-blue-500 hover:text-blue-700 cursor-pointer"
        @click="editHiringNeeds = true; editSalary = false; messageData = {} as BackData"
      >
        {{ hiringNeeds ? hiringNeeds : t(`buildings.hiring.set`) }}
      </span>
      <div v-else class="flex items-center gap-1">
        <InputNumber
          v-model="hiringNeeds"
          class="!w-1/2"
          input-class="!p-2 !w-1/2 !min-w-[60px]"
          show-buttons
        />
        <Button
          icon="pi pi-check"
          severity="info"
          size="small"
          @click="setHiringData(hiringNeeds, 'needs')"
        />
        <Button
          icon="pi pi-times"
          severity="danger"
          size="small"
          @click="editHiringNeeds = false; hiringNeeds = building.hiringNeeds"
        />
      </div>
    </div>
    <div>
      <p>
        {{ t(`buildings.hiring.salary`) }}:
      </p>
      <span
        v-if="!editSalary"
        class="font-bold text-blue-500 hover:text-blue-700 cursor-pointer"
        @click="editSalary = true; editHiringNeeds = false; messageData = {} as BackData"
      >
        {{ salary ? moneyFormat(salary) : t(`buildings.hiring.set`) }}
      </span>
      <div v-else class="flex items-center gap-1">
        <InputNumber
          v-model="salary"
          class="!w-2/3"
          input-class="!p-2 !w-2/3"
          show-buttons
          :step="10"
        />
        <Button
          icon="pi pi-check"
          severity="info"
          size="small"
          @click="setHiringData(salary, 'salary')"
        />
        <Button
          icon="pi pi-times"
          severity="danger"
          size="small"
          @click="editSalary = false; salary = building.salary"
        />
      </div>
    </div>
    <div>
      {{ t(`buildings.coefficient efficiency`) }}:
      {{ (building?.workers / building?.buildingType?.workers).toFixed(2) }}
    </div>
    <hr />
    <div v-if="!isMapFetching">
      {{ t(`buildings.hiring.average salary`) }}:
      {{ moneyFormat(getAverageSalary()) }}
      <Button
        class="w-max"
        :label="t(`buildings.emergency hiring`)"
        severity="secondary"
        size="small"
        @click="confirmEmergencyHiring($event, building._id)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useMutation } from '@tanstack/vue-query'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import { useConfirm } from 'primevue/useconfirm'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { type JsonResult } from '@/api'
import { postBuildingEmergencyHiringMutation } from '@/api/@tanstack/vue-query.gen'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useBuildings } from '@/composables/useBuildings'
import { useMap } from '@/composables/useMap'
import { BackData } from '@/types'
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
const messageData = ref<BackData>()
const confirm = useConfirm()

const { setHiring } = useBuildings()
const { t } = useI18n()
const { getMap } = useMap()
const { data: map, isFetching: isMapFetching } = getMap()
const getAverageSalary = () => {
  if (map.value?.length && props?.building) {
    return  map.value.filter(item => item.x === props.building.x && item.y === props.building.y)[0].averageSalary
  }
  
  return 0
}

const setHiringData = (value, option: HiringOptions) => {
  messageData.value = {} as BackData
  const isSalary = option === 'salary'

  const payload = {
    buildingId: props.building._id,
    salary: isSalary ? value : salary.value,
    hiringNeeds: !isSalary ? value : hiringNeeds.value
  }

  const { data: dataSalary, onFetchResponse: onFetchResponseSalary } = setHiring(payload)

  onFetchResponseSalary(() => {
    messageData.value = dataSalary.value

    editHiringNeeds.value = false
    editSalary.value = false
  })
}

const confirmEmergencyHiring = (event: any, id: string) => {
  confirm.require({
    target: event.currentTarget,
    message: t('common.confirm'),
    icon: 'pi pi-info-circle',
    acceptClass: 'p-button-danger p-button-sm',
    acceptLabel: t('common.yes'),
    rejectLabel: t('common.no'),
    accept: () => handleEmergencyHiring(id)
  })
}

const emergencyHiringMutate = useMutation({
  ...postBuildingEmergencyHiringMutation(),
  onSuccess: (data: JsonResult) => {
    messageData.value = data
    console.log(messageData.value)
  }
})

const handleEmergencyHiring = (id) => {
  const payload = { buildingId: id }
  emergencyHiringMutate.mutate({ body: { ...payload } })
}

</script>
