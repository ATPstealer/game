<template>
  <p class="text-xl mb-4">
    <MessageBlock
      v-if="message"
      :message="message"
    />
  </p>
  <p v-if="resource">
    Sell
    <span class="font-bold">
      {{ resource?.name }}
    </span>
    in
    <span class="font-bold">
      {{ resource?.x }}:{{ resource?.y }}
    </span>
  </p>
  <div v-else>
    <div class="flex flex-col">
      <label class="font-bold text-xl" for="resource-type">Resource type</label>
      <Dropdown
        :options="resourceTypes"
        option-label="name"
        option-value="id"
        placeholder="Choose resource"
        v-model="resourceTypeId"
        input-id="resource-type"
      />
    </div>
    <div class="flex gap-4 mt-4">
      <div class="flex flex-col">
        <label class="font-bold text-xl">X:</label>
        <InputNumber
          v-model="x"
          show-buttons
          :min="-2"
          :max="2"
          input-class="w-12"
        />
      </div>
      <div class="flex flex-col">
        <label class="font-bold text-xl">Y:</label>
        <InputNumber
          v-model="y"
          show-buttons
          :min="-2"
          :max="2"
          input-class="w-12"
        />
      </div>
    </div>
  </div>
  <div class="flex flex-col gap-4 mt-4">
    <div class="flex flex-col">
      <label class="font-bold text-xl">Amount:</label>
      <InputNumber
        v-model="amount"
        show-buttons
        :step="10"
      />
    </div>
    <div class="flex flex-col mb-4">
      <label class="font-bold text-xl">Price per unit:</label>
      <InputNumber
        v-model="priceForUnit"
        show-buttons
      />
    </div>
    <div class="flex gap-4 items-center" v-if="!resource">
      <label class="font-bold text-xl">Order type:</label>
      <SelectButton
        :options="orderTypes"
        option-label="name"
        option-value="value"
        v-model="sell"
        :pt="{
          button: ({ context }) => ({
            class: context.active ? 'bg-blue-400 border-blue-400' : undefined
          })
        }"
      />
    </div>
    <Button
      :label="resource? 'Sell' : 'Create'"
      @click="create"
      class="self-center w-1/2 mt-4"
    />
  </div>
</template>
<script setup lang="ts">
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import InputNumber from 'primevue/inputnumber'
import SelectButton from 'primevue/selectbutton'
import { computed, ref } from 'vue'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useGetData } from '@/composables/useGetData'
import { useOrders } from '@/composables/useOrders'
import { DataMessage } from '@/types'
import { Resource, ResourceType } from '@/types/Resources/index.interface'

interface Props {
  resource?: Resource;
}
const props = defineProps<Props>()
const emits = defineEmits<{(e: 'close'): void}>()

const amount = ref<number>(0)
const priceForUnit = ref<number>(0)
const sell = ref<boolean>(true)
const message = ref<DataMessage | null>(null)
const resourceTypes = ref<ResourceType[]>([])
const resourceTypeId = ref<number>()
const x = ref<number>(0)
const y = ref<number>(0)

const { createOrder } = useOrders()

const orderTypes = computed(() => [
  {
    name: 'Sell',
    value: true
  },
  {
    name: 'Buy',
    value: false
  }
])

if (!props.resource) {
  const { data, onFetchResponse } = useGetData<ResourceType[]>('/resource/types')

  onFetchResponse(() => {
    resourceTypes.value = data.value
  })
}

const create = () => {
  const payload = {
    x: props.resource?.x || x.value,
    y: props.resource?.y || y.value,
    amount: amount.value,
    priceForUnit: priceForUnit.value,
    resourceTypeID: props.resource?.resourceTypeId || resourceTypeId.value,
    sell: sell.value
  }

  const { onFetchResponse, dataMessage } = createOrder(payload)
  onFetchResponse(() => {
    message.value = dataMessage.value

    setTimeout(() => {
      if (message.value?.status === 'success') {
        emits('close')
      }
    }, 2000)
  })
}

</script>
<style scoped>
label {
  @apply mb-2;
}
</style>
