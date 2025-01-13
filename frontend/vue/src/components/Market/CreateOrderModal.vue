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
import { computed, ref, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import {
  type JsonResult,
  type PostMarketOrderCreateMutationRequest,
  type ResourceWithData,
  useGetResourceTypes,
  usePostMarketOrderCreate
} from '@/gen'

interface Props {
  resource?: ResourceWithData | undefined;
}
const props = defineProps<Props>()
const emits = defineEmits<{(e: 'close'): void}>()

const amount = ref<number>(0)
const priceForUnit = ref<number>(0)
const sell = ref<boolean>(true)
const messageData = ref<JsonResult>()
const resourceTypeId = ref<number>()
const x = ref<number>(0)
const y = ref<number>(0)

const { t } = useI18n()

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

const { data: resourceTypesQuery, refetch } = useGetResourceTypes()

const resourceTypes = computed(() => {
  return unref(resourceTypesQuery)?.data || []
})

if (!props.resource) {
  refetch()
}

const createOrderMutate = usePostMarketOrderCreate({
  mutation: {
    onSuccess: data => {
      messageData.value = data

      setTimeout(() => {
        if (!data?.text) {
          emits('close')
        }
      }, 2000)
    }
  }
})

const create = () => {
  messageData.value = {} as JsonResult

  const payload = {
    x: props.resource?.x || x.value,
    y: props.resource?.y || y.value,
    amount: amount.value,
    priceForUnit: priceForUnit.value,
    resourceTypeId: props.resource?.resourceTypeId || resourceTypeId.value,
    sell: sell.value
  }
  // TODO: неверная типизация payload в сваггере
  createOrderMutate.mutate({ data: { ...payload } as PostMarketOrderCreateMutationRequest })
}

</script>
<style scoped>
label {
  @apply mb-2;
}
</style>
