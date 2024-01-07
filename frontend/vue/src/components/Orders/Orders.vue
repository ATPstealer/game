<template>
  <h2 class="mb-5 text-center">
    {{ t(`storages.title`) }}
  </h2>
  <DataTable
    v-if="!isFetching"
    :value="orders"
    size="small"
    striped-rows
  >
    <Column
      :header="t(`map.cell`)"
      class="w-1/5"
    >
      <template #body="{data}: {data: Order}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
    <Column
      field="resourceName"
      :header="t(`orders.columns.resource`)"
      class="w-1/5"
    />
    <Column
      field="amount"
      :header="t(`orders.columns.amount`)"
      class="w-1/5"
    />
    <Column
      field="priceForUnit"
      :header="t(`orders.columns.price`)"
      class="w-1/5"
    >
      <template #body="{data}: {data: Order}">
        {{ moneyFormat(data.priceForUnit) }}
      </template>
    </Column>
    <Column
      :header="t(`orders.columns.type`)"
      class="w-1/5"
    >
      <template #body="{data}: {data: Order}">
        <p class="font-bold text-red-500 hover:text-red-700 block" @click="closeOrder(data.id)">
          {{ data.sell ? 'sell' : 'buy' }}
        </p>
      </template>
    </Column>
  </DataTable>
  <Loading v-else />
</template>
<script setup lang="ts">
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import { useI18n } from 'vue-i18n'
import { moneyFormat } from '../../utils/moneyFormat'
import Loading from '@/components/Common/Loading.vue'
import { useGetData } from '@/composables/useGetData'
import { useOrders } from '@/composables/useOrders'
import { Order, Storage } from '@/types'

const { data: orders, isFetching } = useGetData<Order[]>('/market/order/my')
const { closeOrder } = useOrders()

const { t } = useI18n()
</script>
