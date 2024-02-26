<template>
  <!-- TODO: Сделать перевод; x,y в одну колонку;  -->
  <DataTable
    v-if="orders?.filter(item => item.sell).length"
    :value="orders?.filter(item => item.sell)"
    @row-click="showOrder"
    size="small"
    striped-rows
    :row-class="sellRowClass"
  >
    <Column
      :header="t(`orders.columns.resource`)"
      class="w-1/3"
    >
      <template #body="{data}: {data: Order}">
        {{ t(`resources.types.${data.resourceType.name.toLowerCase()}`) }}
      </template>
    </Column>

    <Column
      field="amount"
      header="Amount"
      class="w-1/6"
    />
    <Column
      field="priceForUnit"
      header="Price"
      class="w-1/6"
    />
    <Column
      field="x"
      header="X"
      class="w-1/6"
    />
    <Column
      field="y"
      header="Y"
      class="w-1/6"
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
      field="amount"
      class="w-1/6"
    >
      <template #header>
        <span class="invisible">Amount</span>
      </template>
    </Column>
    <Column
      field="priceForUnit"
      class="w-1/6"
    >
      <template #header>
        <span class="invisible">Price</span>
      </template>
    </Column>
    <Column
      field="x"
      class="w-1/6"
    >
      <template #header>
        <span class="invisible">X</span>
      </template>
    </Column>
    <Column
      field="y"
      class="w-1/6"
    >
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
    :header="!order.sell ? 'You\'ll sell' : 'You\'ll buy'"
  >
    <div class="flex flex-col gap-4">
      <MessageBlock :message="message" v-if="message" />
      <p><span class="font-bold">{{ order.resourceType.name }}</span></p>
      <p>
        Amount:
        <InputNumber
          input-class="w-28"
          class="h-8"
          show-buttons
          v-model="amount"
        />
      </p>
      <p>Price: <span class="font-bold">{{ order.priceForUnit * amount }}</span></p>
      <Button
        @click="execOrder"
        :label="!order.sell ? 'Sell' : 'Buy'"
        :severity="!order.sell ? 'danger' : 'primary'"
      />
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable, { DataTableRowClickEvent } from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import InputNumber from 'primevue/inputnumber'
import { ref, toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useOrders } from '@/composables/useOrders'
import type { DataMessage, Order } from '@/types'
import type { MarketParams } from '@/types/Resources/index.interface'

interface Props {
  searchParams: MarketParams;
}

const { t } = useI18n()

const props = defineProps<Props>()
const { searchParams } = toRefs(props)

const order = ref<Order>({} as Order)
const showModal = ref<boolean>(false)
const message = ref<DataMessage | null>(null)
const amount = ref<number>(0)

const { getOrders, executeOrder } = useOrders()
const { data: orders, execute } = getOrders(searchParams.value)

const showOrder = (event: DataTableRowClickEvent) => {
  order.value = event.data
  showModal.value = true
}

const buyRowClass = (data) => {
  const index = orders.value.filter(item => !item.sell).findIndex(item => item.id === data.id)
  const bg = index % 2 === 0 ? 'bg-green-50' : 'bg-green-100'

  return [bg, 'cursor-pointer hover:bg-gray-100']
}

const sellRowClass = (data) => {
  const index = orders.value.filter(item => item.sell).findIndex(item => item.id === data.id)
  const bg = index % 2 === 0 ? 'bg-red-50' : 'bg-red-100'

  return [bg, 'cursor-pointer hover:bg-gray-100']
}

const execOrder = () => {
  const { onFetchResponse, dataMessage } = executeOrder({ orderID: order.value.id, amount: amount.value })
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
