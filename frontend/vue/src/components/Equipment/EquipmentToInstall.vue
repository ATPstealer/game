<template>
  <Accordion>
    <AccordionTab :header="`${t('common.add')} ${t('equipment.title').toLocaleLowerCase()}`">
      <MessageBlock
        v-if="messageData?.code"
        class="mb-4"
        v-bind="messageData"
      />
      <DataTable
        :value="availableEquipment"
      >
        <Column :header="t('common.title')">
          <template #body="{data}: {data: Equipment}">
            {{ data.name }}
          </template>
        </Column>
        <Column
          :header="t(`equipment.columns.effect`)"
        >
          <template #body="{data}: {data: Equipment}">
            {{ t(`equipment.effect.${data.effectId.toString()}`) }}
          </template>
        </Column>
        <Column
          :header="t(`common.value`)"
        >
          <template #body="{data}: {data: Equipment}">
            {{ data.value }}
          </template>
        </Column>
        <Column
          :header="t(`map.square`)"
        >
          <template #body="{data}: {data: Equipment}">
            {{ data.square }}
          </template>
        </Column>
        <Column
          :header="t(`equipment.columns.durability`)"
        >
          <template #body="{data}: {data: Equipment}">
            {{ data.durability }}
          </template>
        </Column>
        <Column field="amount" :header="t('common.total')" />
        <Column :header="t('common.amount')">
          <template #body="{data}: {data: Equipment}">
            <InputNumber
              v-model="amounts[data.id]"
              input-class="!w-[100px]"
              :max="data.amount"
              :min="1"
              show-buttons
            />
          </template>
        </Column>
        <Column>
          <template #body="{data}: {data: Equipment}">
            <Button
              :label="t('common.add')"
              size="small"
              @click="add(data)"
            />
          </template>
        </Column>
      </DataTable>
    </AccordionTab>
  </Accordion>
</template>

<script setup lang="ts">
import { useQuery, useMutation } from '@tanstack/vue-query'
import Accordion from 'primevue/accordion'
import AccordionTab from 'primevue/accordiontab'
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import InputNumber from 'primevue/inputnumber'
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { EquipmentType, ResourceAsEquipment, BuildingWithData, JsonResult } from '@/api'
import { getEquipmentMyOptions, postBuildingInstallEquipmentMutation } from '@/api/@tanstack/vue-query.gen'
import MessageBlock from '@/components/Common/MessageBlock.vue'

interface Props {
  building: BuildingWithData;
}

type Equipment = EquipmentType & {amount: number}

const props = defineProps<Props>()

const emits = defineEmits<{(e: 'update-equipment'): void}>()

const myEquipment = ref<ResourceAsEquipment[]>([])
const amounts = ref<Record<number, number>>({})
const messageData = ref<JsonResult>()

const { t } = useI18n()

const { suspense, refetch } = useQuery({
  ...getEquipmentMyOptions({ query: { x: props.building.x, y: props.building.y } }),
  select: (data: any) => {
    myEquipment.value = data.data
  }
})

const addEquipment = useMutation({
  ...postBuildingInstallEquipmentMutation(),
  onSuccess: (data: JsonResult) => {
    messageData.value = data
    refetch()
    emits('update-equipment')
  }
})

const availableEquipment = computed(() => {
  return myEquipment.value.map(item => {
    return {
      amount: item.amount,
      ...item.equipmentType
    }
  }).filter(item => item.amount > 0)
})

const add = (data: Equipment) => {
  const amount = amounts.value[data.id]

  const payload = {
    buildingId: props.building._id,
    equipmentTypeId: data.id,
    amount
  }

  addEquipment.mutate({ body: { ...payload } })
}

const setAmounts = () => {
  amounts.value = availableEquipment.value.reduce((acc: Record<number, number>, cur: EquipmentType) => {
    acc[cur.id] = 1

    return acc
  }, {})
}

watch(availableEquipment, () => {
  if (availableEquipment.value.length) {
    setAmounts()
  }
})

onMounted(() => {
  setAmounts()
})

await suspense()
</script>

<style scoped>

</style>