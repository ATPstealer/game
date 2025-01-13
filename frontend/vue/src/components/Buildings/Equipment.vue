<template>
  <div class="flex flex-col gap-4">
    <Divider />
    <p class="font-bold capitalize">
      {{ t(`common.equipment`) }}
    </p>
    <!--    <Button-->
    <!--      class="w-max"-->
    <!--      label="Добавить"-->
    <!--      size="small"-->
    <!--      @click="openModal('add')"-->
    <!--    />-->
    <!--    <Button-->
    <!--      class="w-max"-->
    <!--      label="Удалить"-->
    <!--      severity="danger"-->
    <!--      size="small"-->
    <!--      @click="openModal('delete')"-->
    <!--    />-->
    <Button
      class="w-max"
      :label="t('equipment.title')"
      size="small"
      @click="router.push({name: 'EquipmentId', params: {id: building?._id}})"
    />
    <Dialog
      v-model:visible="showEquipmentModal"
      dismissable-mask
      header="Доступное оборудование "
      modal
      @hide="hideModal"
    >
      <MessageBlock
        v-if="messageData?.code"
        v-bind="messageData"
      />
      <DataTable
        v-if="process === 'add' && availableEquipment.length"
        :loading="loading"
        :value="availableEquipment"
      >
        <Column header="Название">
          <template #body="{data}: {data: Equipment}">
            {{ data.equipmentType.name }}
          </template>
        </Column>
        <Column
          :header="t(`equipment.columns.effect`)"
        >
          <template #body="{data}: {data: Equipment}">
            {{ t(`equipment.effect.${data.equipmentType.effectId.toString()}`) }}
          </template>
        </Column>
        <Column
          :header="t(`common.value`)"
        >
          <template #body="{data}: {data: Equipment}">
            {{ data.equipmentType.value }}
          </template>
        </Column>
        <Column
          :header="t(`map.square`)"
        >
          <template #body="{data}: {data: Equipment}">
            {{ data.equipmentType.square }}
          </template>
        </Column>
        <Column
          :header="t(`equipment.columns.durability`)"
        >
          <template #body="{data}: {data: Equipment}">
            {{ data.equipmentType.durability }}
          </template>
        </Column>
        <Column field="amount" header="Количество" />
        <Column header="Сколько">
          <template #body="{data}: {data: Equipment}">
            <InputNumber
              v-model="amount[data.equipmentType.id]"
              input-class="!w-[100px]"
              :max="data.amount"
              :min="1"
              show-buttons
            />
          </template>
        </Column>
        <Column header="Добавить">
          <template #body="{data}: {data: Equipment}">
            <Button label="Добавить" @click="changeEquipment(data, true)" />
          </template>
        </Column>
      </DataTable>
      <DataTable
        v-if="process === 'delete' && availableEquipment?.length"
        :loading="loading"
        :value="availableEquipment"
      >
        <Column field="name" header="Название" />
        <Column field="amount" header="Количество" />
        <Column field="durability" header="Прочность" />
        <Column header="Сколько">
          <template #body="{data}: {data: EquipmentType & {amount: number}}">
            <InputNumber
              v-model="amount[data.id]"
              input-class="!w-[100px]"
              :max="data.amount"
              :min="1"
              show-buttons
            />
          </template>
        </Column>
        <Column header="Удалить">
          <template #body="{data}: {data: Equipment}">
            <Button
              label="Удалить"
              severity="danger"
              @click="changeEquipment(data, false)"
            />
          </template>
        </Column>
      </DataTable>
      <div v-if="!loading && !availableEquipment?.length">
        Доступное оборудование отсутствует
      </div>
    </Dialog>

    <div v-if="building?.equipmentEffect?.length && blueprints?.length">
      <p class="font-bold mb-4">
        Эффекты оборудования:
      </p>
      <div class="flex flex-col gap-2">
        <Tag
          v-for="effect in building.equipmentEffect"
          :key="effect.effectId"
          v-tooltip.top="`Влияние: ${effect.value.toFixed(3)} на ${effect.blueprintId === 0 ? 'все чертежи' : blueprints.find(bp => bp.id === effect.blueprintId)?.name}`"
          severity="info"
          :value="t(`equipment.effect.${effect.effectId.toString()}`)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import Divider from 'primevue/divider'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useBuildings } from '@/composables/useBuildings'
import { useGetData } from '@/composables/useGetData'
import type { BuildingWithData } from '@/gen'
import router from '@/router'
import type { BackData } from '@/types'
import type { Blueprint } from '@/types/Buildings/index.interface'
import type { Equipment, EquipmentType } from '@/types/Equipment/index.interface'

interface Props {
  building: BuildingWithData | undefined;
}

const props = defineProps<Props>()
const execute = inject<() => void>('execute')

const showEquipmentModal = ref<boolean>(false)
const availableEquipment = ref<any[]>([])
const loading = ref<boolean>(false)
const amount = ref<any>()
const process = ref<string>('')
const messageData = ref<BackData>()

const { t } = useI18n()
const { installEquipment } = useBuildings()
const { data: equipmentTypes, onFetchResponse } = useGetData<EquipmentType[]>('/equipment/types')
const { data: blueprints } = useGetData<Blueprint[]>('/building/blueprints')
const { data: currentEquipment, onFetchResponse: onCurrentEquipmentResponse, execute: getCurrentEquipment } = useGetData<Equipment[]>(`/equipment/my?x=${props.building?.x}&y=${props.building?.y}`, false)

const openModal = (event: string) => {
  process.value = event
  loading.value = true
  showEquipmentModal.value = true

  if (event === 'add') {
    getCurrentEquipment()
    onCurrentEquipmentResponse(() => {
      availableEquipment.value = [...currentEquipment.value]

      amount.value = availableEquipment.value.reduce((acc: Record<string, number>, cur: Equipment) => {
        acc[cur.equipmentType.id] = 1

        return acc
      }, {})
      console.log(availableEquipment.value)
      loading.value = false
    })
  }

  else {
    if (!props?.building?.equipment?.length) {
      availableEquipment.value = []
      loading.value = false

      return
    }
    const localEquipment = props.building.equipment.map(eq => eq.equipmentTypeId)

    availableEquipment.value = [...equipmentTypes.value.filter(eq => localEquipment.includes(eq.id))].map(eq => {
      return {
        ...eq,
        amount: props.building?.equipment?.find(item => item.equipmentTypeId === eq.id)?.amount
      }
    })

    amount.value = availableEquipment.value.reduce((acc: Record<string, number>, cur: EquipmentType) => {
      acc[cur.id] = 1

      return acc
    }, {})

    loading.value = false
  }
}

const changeEquipment = (equipment: any, add: boolean) => {
  interface Payload {
    buildingId: string;
    equipmentTypeId?: number;
    amount?: number;
  }

  const payload: Payload = {
    buildingId: props.building!._id
  }

  if (add) {
    payload.equipmentTypeId = equipment.equipmentType.id
    payload.amount = amount.value[equipment.equipmentType.id]
  } else {
    payload.equipmentTypeId = equipment.id
    payload.amount = -amount.value[equipment.id]
  }

  const { data, onFetchResponse } = installEquipment(payload)

  onFetchResponse( () => {
    if (add) {
      amount.value[equipment.equipmentType.id] = 1
    }
    else {
      amount.value[equipment.id] = 1
    }
    messageData.value = data.value

    setTimeout( () => {
      if (messageData.value?.code === 34) {
        messageData.value= {} as BackData
      }
      else {
        execute?.()
      }
    }, 1000)
  })
}

const hideModal = async () => {
  process.value = ''
  showEquipmentModal.value = false
  availableEquipment.value = []
  messageData.value = {} as BackData
}

</script>