<template>
  <DataTable
    v-if="orders?.filter(item => item.sell).length"
    :value="orders?.filter(item => item.sell)"
    @row-click="showOrder"
    size="small"
    striped-rows
    :row-class="sellRowClass"
  >
    <Column
      field="resourceName"
      header="Name"
    />
    <Column
      field="amount"
      header="Amount"
    />
    <Column
      field="priceForUnit"
      header="Price"
    />
    <Column
      field="x"
      header="X"
    />
    <Column
      field="y"
      header="Y"
    />
  </DataTable>
  <DataTable
    v-if="orders?.filter(item => !item.sell).length"
    :value="orders?.filter(item => !item.sell)"
    @row-click="showOrder"
    size="small"
    striped-rows
    :row-class="buyRowClass"
    :pt="{
      headerRow: {
        class: 'bg-white'
      },
    }"
  >
    <Column field="resourceName">
      <template #header>
        <span class="invisible">Name</span>
      </template>
    </Column>
    <Column field="amount">
      <template #header>
        <span class="invisible">Amount</span>
      </template>
    </Column>
    <Column field="priceForUnit">
      <template #header>
        <span class="invisible">Price</span>
      </template>
    </Column>
    <Column field="x">
      <template #header>
        <span class="invisible">X</span>
      </template>
    </Column>
    <Column field="y">
      <template #header>
        <span class="invisible">Y</span>
      </template>
    </Column>
  </DataTable>
  <div v-if="!orders?.length">
    Этого ресурса на рынке нет
  </div>
  <Dialog
    v-model:visible="showModal"
    :style="{ width: '25rem' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :header="order.sell ? 'You\'ll sell' : 'You\'ll buy'"
  >
    <div class="flex flex-col gap-4">
      <MessageBlock :message="message" v-if="message" />
      <p><span class="font-bold">{{ order.resourceName }}</span></p>
      <p>Amount: <span class="font-bold">{{ order.amount }}</span></p>
      <p>Price: <span class="font-bold">{{ order.priceForUnit * order.amount }}</span></p>
      <Button
        @click="execOrder"
        :label="order.sell ? 'Sell' : 'Buy'"
        :severity="order.sell ? 'danger' : 'primary'"
      />
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable, { DataTableRowClickEvent } from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import { ref, toRefs } from 'vue'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useOrders } from '@/composables/useOrders'
import type { DataMessage, Order } from '@/types'
import type { MarketParams } from '@/types/Resources/index.interface'

interface Props {
 searchParams: MarketParams;
}

const props = defineProps<Props>()
const { searchParams } = toRefs(props)

const order = ref<Order>({} as Order)
const showModal = ref<boolean>(false)
const message = ref<DataMessage | null>(null)

const { getOrders, executeOrder } = useOrders()
const { data: orders, execute } = getOrders(searchParams.value)

const showOrder = (event: DataTableRowClickEvent) => {
  order.value = event.data
  showModal.value = true
}

const buyRowClass = (data) => {
  const index = orders.value.filter(item => !item.sell).findIndex(item => item.id === data.id)
  const bg =  index % 2 === 0 ? 'bg-green-50' : 'bg-green-100'

  return [bg, 'cursor-pointer hover:bg-gray-100']
}

const sellRowClass = (data) => {
  const index = orders.value.filter(item => item.sell).findIndex(item => item.id === data.id)
  const bg =  index % 2 === 0 ? 'bg-red-50' : 'bg-red-100'

  return [bg, 'cursor-pointer hover:bg-gray-100']
}

const execOrder = () => {
  const { onFetchResponse, dataMessage } = executeOrder(order.value.id)
  onFetchResponse(() => {
    message.value = dataMessage.value

    setTimeout(() => {
      showModal.value = false
      message.value = null
      execute()
    }, 2000)
  })
}
</script>

<style scoped>

</style>
