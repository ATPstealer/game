<template>
  <DataTable
    v-if="tableData?.length"
    edit-mode="cell"
    :pt="{
      bodyRow: {
        class: 'h-14'
      }
    }"
    size="small"
    striped-rows
    :value="tableData"
    @cell-edit-complete="onCellEditComplete"
  >
    <Column
      field="name"
      :header="t(`buildings.store.columns.name`)"
    >
      <template #body="{data}:{data: ResourceType}">
        {{ t(`resources.types.${data.name.toLowerCase()}`) }}
      </template>
    </Column>
    <Column
      class="min-w-[150px] max-w-[150px]"
      field="price"
      :header="t(`buildings.store.columns.price`)"
    >
      <template #body="{ data, field }">
        <span class="col-span-1 font-bold text-blue-500 hover:text-blue-700 block">
          {{ data[field] ? moneyFormat(data[field]) : t('buildings.store.setPrice') }}
        </span>
      </template>
      <template #editor="{ data, field }">
        <InputNumber
          v-model="data[field]"
          autofocus
          class="!w-1/2"
          input-class="!p-2 !w-1/2"
          max-fraction-digits="2"
          min-fraction-digits="2"
        />
      </template>
    </Column>
    <Column
      field="sellSum"
      :header="t(`buildings.store.columns.sell count`)"
    />
    <Column
      field="revenue"
      :header="t(`buildings.store.columns.revenue today`)"
    />
    <Column
      :header="t(`buildings.store.columns.status`)"
    >
      <template #body="{data}: {data: Goods}">
        {{ t(`buildings.store.status.${data.status.toLowerCase()}`) }}
      </template>
    </Column>
  </DataTable>
</template>
<script setup lang="ts">
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import InputNumber from 'primevue/inputnumber'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBuildings } from '@/composables/useBuildings'
import { useGetData } from '@/composables/useGetData'
import { Building, Goods } from '@/types/Buildings/index.interface'
import { ResourceType } from '@/types/Resources/index.interface'
import { moneyFormat } from '@/utils/moneyFormat'

interface Props {
  building: Building;
}

const props = defineProps<Props>()
const resourcesTypes = ref<ResourceType[]>([])
const updatedBuilding = ref<Building>(props.building)

const { data, onFetchResponse, isFetching: isFetchingResourcesTypes } = useGetData<ResourceType[]>('/resource/types')
onFetchResponse(() => {
  resourcesTypes.value = data.value.filter(item => item.storeGroup === props.building.buildingType.buildingSubGroup)
})

const executeBuildingUpdate = () => {
  const { data: building, onFetchResponse } = useGetData<Building[]>(`/building/my?_id=${  props.building._id}`)
  onFetchResponse(() => {
    updatedBuilding.value = building.value[0]
  })
}

const tableData = computed(() => {
  return resourcesTypes.value.map(item => {
    return {
      name: item.name,
      price: getGoodsData(item.id)?.price || 0,
      resourceTypeId: item.id,
      sellSum: getGoodsData(item.id)?.sellSum || 0,
      revenue: getGoodsData(item.id)?.revenue || 0,
      status: getGoodsData(item.id)?.status || 'notSet'
    }
  })
})

const getGoodsData = (id: number) => {
  if (updatedBuilding.value.goods) {
    return updatedBuilding.value.goods.find(item => item.resourceTypeId === id)
  }

  return null
}

const { setPrice } = useBuildings()
const onCellEditComplete = (event) => {
  const payload = {
    buildingId: props.building._id,
    resourceTypeId: event.data.resourceTypeId,
    price: event.newValue
  }
  const { onFetchResponse } = setPrice(payload)

  onFetchResponse(() => {
    executeBuildingUpdate()
  })
}

const { t } = useI18n()
</script>
