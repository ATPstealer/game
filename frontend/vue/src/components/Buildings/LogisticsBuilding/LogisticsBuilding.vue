<template>
  <h3>{{ t('logistics.title') }}</h3>
  <Card class="mt-8">
    <template #content>
      <MessageBlock v-if="messageData" v-bind="messageData" />
      <p class="mb-2">
        <span class="font-bold">{{ t('logistics.capacity') }}</span>: {{ data?.capacity }}
      </p>
      <p class="mb-2">
        <span class="font-bold">Max.</span> <span class="font-bold lowercase">{{ t('logistics.capacity') }}</span>: {{ data?.capacityMax }}
      </p>
      <p class="mb-2">
        <span class="font-bold">{{ t('logistics.speed') }}</span>: {{ data?.speed }}
      </p>
      <div class="flex flex-col gap-2">
        <span class="font-bold">{{ t('common.price') }}</span>
        <div class="flex gap-4">
          <InputNumber
            v-model="price"
            input-class="max-w-[150px]"
            show-buttons
          />
          <Button
            class="w-max"
            :label="t('buildings.hiring.set')"
            @click="setPrice"
          />
        </div>
      </div>
    </template>
  </Card>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Card from 'primevue/card'
import InputNumber from 'primevue/inputnumber'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { type BuildingWithData, type JsonResult, type Logistics, usePostLogisticsSetPrice } from '@/gen'

interface Props {
  building: BuildingWithData;
}

const props = defineProps<Props>()

const mutationLogisticSetPrice = usePostLogisticsSetPrice({
  mutation: {
    onSuccess: data => {
      messageData.value = data
    }
  }
})

const data = ref<Logistics>(props.building.logistics || {} as Logistics)

const price = ref<number>(data.value.price || 0)
const oldPrice = ref<number>(data.value.price || 0)
const messageData = ref<JsonResult>()

const { t } = useI18n()

const setPrice = () => {
  if (price.value === oldPrice.value) {
    return
  }
  oldPrice.value = price.value

  const payload = {
    buildingId: props.building._id,
    price: price.value
  }

  mutationLogisticSetPrice.mutate({ data: { ...payload } })
}

</script>

<style scoped>

</style>