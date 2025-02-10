<template>
  <h2 class="mb-5 text-center">
    {{ t(`orders.title`) }}
  </h2>
  <DataTable
    v-if="!isFetching"
    size="small"
    striped-rows
    :value="orders"
  >
    <Column
      class="w-1/5"
      :header="t(`map.cell`)"
    >
      <template #body="{data}: {data: OrderWithData}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
    <Column
      class="w-1/5"
      field="Resource"
      :header="t(`orders.columns.resource`)"
    >
      <template #body="{data}: {data: OrderWithData}">
        {{ t(`resources.types.${data.resourceType.name.toLowerCase()}`) }}
      </template>
    </Column>
    <Column
      class="w-1/5"
      field="amount"
      :header="t(`orders.columns.amount`)"
    />
    <Column
      class="w-1/5"
      field="priceForUnit"
      :header="t(`orders.columns.price`)"
    >
      <template #body="{data}: {data: OrderWithData}">
        {{ moneyFormat(data.priceForUnit) }}
      </template>
    </Column>
    <Column
      class="w-1/5"
      :header="t(`orders.columns.type`)"
    >
      <template #body="{data}: {data: OrderWithData}">
        <p class="font-bold text-red-500 hover:text-red-700 block" @click="orderMutation.mutate({params: {orderId: data._id}})">
          {{ data.sell ? 'sell' : 'buy' }}
        </p>
      </template>
    </Column>
  </DataTable>
  <Loading v-else />
</template>
<script setup lang="ts">
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import { computed, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import Loading from '@/components/Common/Loading.vue'
import { type OrderWithData, useDeleteMarketOrderClose, useGetMarketOrderMy } from '@/gen'
import { moneyFormat } from '@/utils/moneyFormat'

const { t } = useI18n()

const { data: ordersQuery, isFetching, suspense } = useGetMarketOrderMy()
await suspense()
const orders = computed(() => unref(ordersQuery)?.data)

const orderMutation = useDeleteMarketOrderClose()

</script>
