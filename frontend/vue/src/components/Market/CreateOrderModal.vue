<template>
  <p class="text-xl mb-4">
    <MessageBlock
      v-if="messageData?.code"
      v-bind="messageData"
    />
  </p>
  <p v-if="resource">
    {{ t('common.sell') }}
    <span class="font-bold">
      {{ t(`resources.types.${resource?.resourceType.name.toLowerCase()}`) }}
    </span>
    {{ t('common.in') }}
    <span class="font-bold">
      {{ resource?.x }}:{{ resource?.y }}
    </span>
  </p>
  <div v-else>
    <div class="flex flex-col">
      <label class="font-bold text-xl" for="resource-type">{{ t('resources.resourceType') }}</label>
      <Dropdown
        v-model="resourceTypeId"
        input-id="resource-type"
        :option-label="event => t(`resources.types.${event.name.toLowerCase()}`)"
        option-value="id"
        :options="resourceTypes"
        :placeholder="t('resources.choose')"
      />
    </div>
    <div class="flex gap-4 mt-4">
      <div class="flex flex-col">
        <label class="font-bold text-xl">X:</label>
        <InputNumber
          v-model="x"
          input-class="w-12"
          :max="2"
          :min="-2"
          show-buttons
        />
      </div>
      <div class="flex flex-col">
        <label class="font-bold text-xl">Y:</label>
        <InputNumber
          v-model="y"
          input-class="w-12"
          :max="2"
          :min="-2"
          show-buttons
        />
      </div>
    </div>
  </div>
  <div class="flex flex-col gap-4 mt-4">
    <div class="flex flex-col">
      <label class="font-bold text-xl">{{ t('common.amount') }}:</label>
      <InputNumber
        v-model="amount"
        show-buttons
        :step="10"
      />
    </div>
    <div class="flex flex-col mb-4">
      <label class="font-bold text-xl">{{ t('resources.sell.price') }}:</label>
      <InputNumber
        v-model="priceForUnit"
        show-buttons
      />
    </div>
    <div v-if="!resource" class="flex gap-4 items-center">
      <label class="font-bold text-xl">{{ t('orders.create.type') }}</label>
      <SelectButton
        v-model="sell"
        option-label="name"
        option-value="value"
        :options="orderTypes"
        :pt="{
          button: ({ context }) => ({
            class: context.active ? 'bg-blue-400 border-blue-400' : undefined
          })
        }"
      />
    </div>
    <Button
      class="self-center w-1/2 mt-4"
      :disabled="!amount || !priceForUnit || (!resource && !resourceTypeId)"
      :label="resource? t('common.sell') : t('common.create')"
      @click="create"
    />
  </div>
</template>
<script setup lang="ts">
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import InputNumber from 'primevue/inputnumber'
import SelectButton from 'primevue/selectbutton'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useGetData } from '@/composables/useGetData'
import { useOrders } from '@/composables/useOrders'
import { BackData } from '@/types'
import { Resource, ResourceType } from '@/types/Resources/index.interface'

interface Props {
  resource?: Resource;
}
const props = defineProps<Props>()
const emits = defineEmits<{(e: 'close'): void}>()

const amount = ref<number>(0)
const priceForUnit = ref<number>(0)
const sell = ref<boolean>(true)
const messageData = ref<BackData>()
const resourceTypes = ref<ResourceType[]>([])
const resourceTypeId = ref<number>()
const x = ref<number>(0)
const y = ref<number>(0)

const { t } = useI18n()
const { createOrder } = useOrders()

const orderTypes = computed(() => [
  {
    name: t('common.sell'),
    value: true
  },
  {
    name: t('common.buy'),
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
  messageData.value = {} as BackData

  const payload = {
    x: props.resource?.x || x.value,
    y: props.resource?.y || y.value,
    amount: amount.value,
    priceForUnit: priceForUnit.value,
    resourceTypeID: props.resource?.resourceTypeId || resourceTypeId.value,
    sell: sell.value
  }

  const { data, onFetchResponse } = createOrder(payload)
  onFetchResponse(() => {
    messageData.value = data.value

    setTimeout(() => {
      if (!data.value.text) {
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
