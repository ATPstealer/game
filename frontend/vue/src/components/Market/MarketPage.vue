<template>
  <DataTable
    v-if="orders?.filter(item => item.sell).length"
    :row-class="sellRowClass"
    size="small"
    striped-rows
    :value="orders?.filter(item => item.sell)"
    @row-click="showOrder"
  >
    <Column
      class="w-1/3"
      :header="t(`orders.columns.resource`)"
    >
      <template #body="{data}: {data: Order}">
        {{ t(`resources.types.${data.resourceType.name.toLowerCase()}`) }}
      </template>
    </Column>

    <Column
      class="w-1/6"
      field="amount"
      :header="t('orders.columns.amount')"
    />
    <Column
      class="w-1/6"
      field="priceForUnit"
      :header="t('orders.columns.price')"
    />
    <Column
      class="w-1/6"
      :header="t(`map.cell`)"
    >
      <template #body="{data}: {data: Order}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
  </DataTable>
  <DataTable
    v-if="orders?.filter(item => !item.sell).length"
    :pt="{
      headerRow: {
        class: 'bg-white'
      },
    }"
    :row-class="buyRowClass"
    size="small"
    striped-rows
    :value="orders?.filter(item => !item.sell)"
    @row-click="showOrder"
  >
    <Column
      class="w-1/3"
    >
      <template #header>
        <span class="invisible">Amount</span>
      </template>
      <template #body="{data}: {data: Order}">
        {{ t(`resources.types.${data.resourceType.name.toLowerCase()}`) }}
      </template>
    </Column>

    <Column
      class="w-1/6"
      field="amount"
    >
      <template #header>
        <span class="invisible">{{ t('orders.columns.amount') }}</span>
      </template>
    </Column>
    <Column
      class="w-1/6"
      field="priceForUnit"
    >
      <template #header>
        <span class="invisible">{{ t('orders.columns.price') }}</span>
      </template>
    </Column>
    <Column
      class="w-1/6"
    >
      <template #header>
        <span class="invisible">{{ t('map.cell') }}</span>
      </template>
      <template #body="{data}: {data: Order}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
  </DataTable>
  <div v-if="!orders?.length">
    Этого ресурса на рынке нет
  </div>
  <Dialog
    v-model:visible="showModal"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :header="!order.sell ? 'You\'ll sell' : 'You\'ll buy'"
    :style="{ width: '25rem' }"
    @hide="messageData = {} as BackData; amount = 0"
  >
    <div class="flex flex-col gap-4">
      <MessageBlock
        v-if="messageData?.code"
        v-bind="messageData"
      />
      <p><span class="font-bold">{{ order.resourceType.name }}</span></p>
      <p>
        Amount:
        <InputNumber
          v-model="amount"
          class="h-8"
          input-class="w-28"
          show-buttons
        />
      </p>
      <p>Price: <span class="font-bold">{{ order.priceForUnit * amount }}</span></p>
      <Button
        :label="!order.sell ? 'Sell' : 'Buy'"
        :severity="!order.sell ? 'danger' : 'primary'"
        @click="execOrder"
      />
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable, { type DataTableRowClickEvent } from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import InputNumber from 'primevue/inputnumber'
import { ref, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useOrders } from '@/composables/useOrders'
import { useGetOrders } from '@/gen'
import type { BackData, Order } from '@/types'
import type { MarketParams } from '@/types/Resources/index.interface'

interface Props {
  searchParams: MarketParams;
}

const { t } = useI18n()

const props = defineProps<Props>()
const { searchParams } = toRefs(props)

const order = ref<Order>({} as Order)
const showModal = ref<boolean>(false)
const messageData = ref<BackData>()
const amount = ref<number>(0)

const { getOrders, executeOrder } = useOrders()
const { data: orders, execute } = getOrders(searchParams.value)

const { data: orders2, refetch: execute2, suspense } = useGetOrders(searchParams.value)
await suspense()

const showOrder = (event: DataTableRowClickEvent) => {
  order.value = event.data
  showModal.value = true
}

const buyRowClass = (data) => {
  const index = orders.value.filter(item => !item.sell).findIndex(item => item._id === data._id)
  const bg = index % 2 === 0 ? 'bg-green-50' : 'bg-green-100'

  return [bg, 'cursor-pointer hover:bg-gray-100']
}

const sellRowClass = (data) => {
  const index = orders.value.filter(item => item.sell).findIndex(item => item._id === data._id)
  const bg = index % 2 === 0 ? 'bg-red-50' : 'bg-red-100'

  return [bg, 'cursor-pointer hover:bg-gray-100']
}

const execOrder = () => {
  messageData.value = {} as BackData

  const { data, onFetchResponse } = executeOrder({ orderID: order.value._id, amount: amount.value })
  onFetchResponse(() => {
    messageData.value = data.value

    setTimeout(() => {
      if (!data.value.text) {
        showModal.value = false
        messageData.value = {} as BackData
        execute()
      }
    }, 2000)
  })
}

</script>

<style scoped>

</style>
