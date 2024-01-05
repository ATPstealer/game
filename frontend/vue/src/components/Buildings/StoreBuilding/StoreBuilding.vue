<template>
  <!-- TODO: пофиксить t() -->
  <h3 class="mb-5 text-center">
    {{ t(`buildings.store.types.${building.buildingSubGroup?.toLowerCase()}`) }}
    {{ t(`buildings.store.name`) }}
  </h3>

  <DataTable
      v-if="tableData?.length"
      :value="tableData"
      size="small"
      striped-rows
      editMode="cell"
      @cell-edit-complete="onCellEditComplete"
  >
    <Column
        field="name"
        :header="t(`buildings.store.columns.name`)"
    />
    <Column
        field="price"
        :header="t(`buildings.store.columns.price`)"
    >
      <template #body="{ data, field }">
        <span class="col-span-1 font-bold text-blue-500 hover:text-blue-700">
          {{ moneyFormat(data[field]) }}
        </span>
      </template>
      <!--  TODO: сделать чтобы не разъезжалось и были ровыне поля при редактировании цены -->>
      <template #editor="{ data, field }">
        <InputNumber v-model="data[field]" autofocus/>
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
        field="status"
        :header="t(`buildings.store.columns.status`)"
    />

  </DataTable>
</template>
<script setup lang="ts">
import {useI18n} from 'vue-i18n'
import {Building, Goods} from "@/types/Buildings/index.interface";
import {useGetData} from "@/composables/useGetData";
import {ResourceType} from "@/types/Resources/index.interface";
import {computed, ref} from 'vue'
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import InputNumber from "primevue/inputnumber";
import {useBuildings} from "@/composables/useBuildings";
import {moneyFormat} from "@/utils/moneyFormat";

interface Props {
  building: Building;
}

const props = defineProps<Props>()
const resourcesTypes = ref<ResourceType[]>([])
const goods = ref<Goods[]>([])

const {data, onFetchResponse, isFetching: isFetchingResourcesTypes} = useGetData('/resource/types')
onFetchResponse(() => {
  resourcesTypes.value = data.value.filter(item => item.storeGroup === props.building.buildingSubGroup)
})

const {data: goodsData, onFetchResponse: onGoodsResponse, execute: executeGoods} = useGetData('/store/goods/get?building_id=' + props.building.id)
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
      status: getGoodsData(item.id)?.status || "Price not set",
    }
  })
})
const getGoodsData = (id: number) => {
  return goods.value.find(item => item.resourceTypeId === id)
}

const {setPrice} = useBuildings()
const onCellEditComplete = async (event) => {
  const payload = {
    buildingId: props.building.id,
    resourceTypeId: event.data.resourceTypeId,
    price: event.newValue,
  }
  await setPrice(payload)
  // TODO: Илья пофиксь плз. Нужно чтобы executeGoods() вызывалось после setPrice()
  setTimeout(() => {
    executeGoods()
  }, 1000)

}

const {t} = useI18n()
</script>
