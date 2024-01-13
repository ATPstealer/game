<template>
  <MessageBlock
    v-if="message"
    :message="message"
    class="mb-4"
  />
  <div class="m-0">
    <div class="font-bold">
      {{ t(`buildings.hiring.name`) }}
    </div>
    <label class="font-bold"> {{ t(`buildings.hiring.workers`) }} / {{ t(`buildings.hiring.hiring needs`) }} </label>
    <div>
      {{ props.building.workers }} /
      <span
        v-if="!editHiringNeeds"
        @click="editHiringNeeds = true; message = null"
        class="font-bold text-blue-500 hover:text-blue-700"
      >
        {{ hiringNeeds ? hiringNeeds : t(`buildings.hiring.set`) }}
      </span>
      <InputNumber
        v-else
        v-model="hiringNeeds"
        show-buttons
        @blur="setHiringNeeds"
        :step="1"
        input-class="!p-2 !w-1/2"
        class="!w-1/2"
      />
    </div>
    <div>
      <label class="font-bold"> {{ t(`buildings.hiring.salary`) }} </label>
    </div>
    <div
      v-if="!editSalary"
      @click="editSalary = true; message = null"
      class="font-bold text-blue-500 hover:text-blue-700"
    >
      {{ salary ? moneyFormat(salary) : t(`buildings.hiring.set`) }}
    </div>
    <InputNumber
      v-else
      v-model="salary"
      show-buttons
      @blur="setSalary"
      :step="10"
      input-class="!p-2 !w-2/3"
      class="!w-2/3"
    />
  </div>
</template>

<script setup lang="ts">
import InputNumber from 'primevue/inputnumber'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useBuildings } from '@/composables/useBuildings'
import { Message } from '@/types'
import { Building } from '@/types/Buildings/index.interface'
import { moneyFormat } from '@/utils/moneyFormat'

interface Props {
  building: Building;
}

const editSalary = ref<boolean>(false)
const editHiringNeeds = ref<boolean>(false)
const props = defineProps<Props>()
const salary = ref<number>(props.building.salary)
const hiringNeeds = ref<number>(props.building.hiringNeeds)
const { setHiring } = useBuildings()
const message = ref<Message | null>(null)
const { t } = useI18n()

const setSalary = (event) => {
  const payload = {
    buildingId: props.building.id,
    salary: Number(event.value),
    hiringNeeds: hiringNeeds.value
  }
  const { dataMessage: dataMessageSalary, onFetchResponse: onFetchResponseSalary } = setHiring(payload)
  onFetchResponseSalary(() => {
    message.value = dataMessageSalary.value
  })
  editSalary.value = false
}
const setHiringNeeds = (event) => {
  const payload = {
    buildingId: props.building.id,
    salary: salary.value,
    hiringNeeds: Number(event.value)
  }
  const { dataMessage: dataMessageHiringNeeds, onFetchResponse: onFetchResponseHiringNeeds } = setHiring(payload)
  onFetchResponseHiringNeeds(() => {
    message.value = dataMessageHiringNeeds.value
  })
  editHiringNeeds.value = false
}

</script>
