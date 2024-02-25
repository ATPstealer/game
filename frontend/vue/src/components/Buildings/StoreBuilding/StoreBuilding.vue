<template>
  <h3 class="mb-5 text-center" v-if="building && building.buildingType && building.buildingType.buildingSubGroup">
    {{ t(`buildings.store.types.${building.buildingType.buildingSubGroup.toLowerCase()}`) }}
    {{ t(`buildings.store.name`) }}
  </h3>

  <DataTable
    v-if="tableData?.length"
    :value="tableData"
    size="small"
    striped-rows
    edit-mode="cell"
    @cell-edit-complete="onCellEditComplete"
    :pt="{
      bodyRow: {
        class: 'h-14'
      }
    }"
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
      field="price"
      :header="t(`buildings.store.columns.price`)"
      class="min-w-[150px] max-w-[150px]"
    >
      <template #body="{ data, field }">
        <span class="col-span-1 font-bold text-blue-500 hover:text-blue-700 block">
          {{ data[field] ? moneyFormat(data[field]) : t('buildings.store.setPrice') }}
        </span>
      </template>
      <!--  TODO: сделать чтобы не разъезжалось и были ровыне поля при редактировании цены -->>
      <template #editor="{ data, field }">
        <InputNumber
          v-model="data[field]"
          autofocus
          input-class="!p-2 !w-1/2"
          class="!w-1/2"
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
const goods = ref<Goods[]>([])

const { data, onFetchResponse, isFetching: isFetchingResourcesTypes } = useGetData<ResourceType[]>('/resource/types')
onFetchResponse(() => {
  resourcesTypes.value = data.value.filter(item => item.storeGroup === props.building.buildingType.buildingSubGroup)
})

const { data: goodsData, onFetchResponse: onGoodsResponse, execute: executeGoods } = useGetData<Goods[]>(`/store/goods/get?building_id=${ props.building._id}`)
onGoodsResponse(() => {
  goods.value = goodsData.value
})

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
  if (goods.value) {
    return goods.value.find(item => item.resourceTypeId === id)
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
    executeGoods()
  })
}

const { t } = useI18n()
</script>
