<template>
  <div class="w-full">
    <MessageBlock
      v-if="messageData?.code"
      class="mb-4"
      v-bind="messageData"
    />
    <DataTable
      v-if="availableEquipment.length"
      :value="availableEquipment"
    >
      <Column field="name" :header="t('common.title')" />
      <Column field="amount" :header="t('common.total')" />
      <Column field="durability" :header="t('equipment.columns.durability')" />
      <Column :header="t('common.amount')">
        <template #body="{data}: {data: EquipmentType & {amount: number}}">
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
        <template #body="{data}: {data: EquipmentType}">
          <Button
            :label="t('common.delete')"
            severity="danger"
            size="small"
            @click="changeEquipment(data)"
          />
        </template>
      </Column>
    </DataTable>
    <p v-else>
      {{ t('equipment.empty') }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { useMutation } from '@tanstack/vue-query'
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import InputNumber from 'primevue/inputnumber'
import { onMounted, ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { BuildingWithData, EquipmentType, JsonResult } from '@/api'
import { postBuildingInstallEquipmentMutation } from '@/api/@tanstack/vue-query.gen'
import MessageBlock from '@/components/Common/MessageBlock.vue'

interface Props {
  building: BuildingWithData;
  equipmentTypes: EquipmentType[];
}

const props = defineProps<Props>()

const emits = defineEmits<{(e: 'update-equipment'): void}>()

const amounts = ref<Record<number, number>>({})
const messageData = ref<JsonResult>()

const availableEquipment = computed(() => {
  const allEquipment = props.building.equipment

  if (!allEquipment?.length) {
    return []
  }

  return props.equipmentTypes.filter(eq => {
    if (allEquipment.find(item => item.equipmentTypeId === eq.id) !== undefined) {
      return eq
    }
  }).map(eq => {
    return {
      ...eq,
      amount: allEquipment.find(item => item.equipmentTypeId === eq.id)?.amount
    }
  })
})

const { t } = useI18n()

const updateEquipment = useMutation({
  ...postBuildingInstallEquipmentMutation(),
  onSuccess: (data: JsonResult) => {
    messageData.value = data
    emits('update-equipment')
  }
})

const changeEquipment = (data: EquipmentType) => {
  const amount = amounts.value[data.id]

  const payload = {
    buildingId: props.building._id,
    equipmentTypeId: data.id,
    amount: -amount
  }

  updateEquipment.mutate({ body: { ...payload } })
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
</script>

<style scoped>

</style>