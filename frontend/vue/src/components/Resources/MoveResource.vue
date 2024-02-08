<template>
  <div class="flex flex-col gap-4">
    <MessageBlock
      v-if="message"
      :message="message"
    />
    <p class="text-xl">
      {{ t('common.move') }}
      <span class="font-bold">
        {{ t(`resources.types.${resource.name.toLowerCase()}`) }}
      </span>
      {{ t('common.from') }}
      <span class="font-bold">
        {{ resource.x }}:{{ resource.y }}
      </span>
      {{ t('common.to') }}
    </p>
    <div class="flex gap-4">
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
    <div class="flex flex-col">
      <label class="font-bold text-xl">{{ t('common.amount') }}:</label>
      <InputNumber
        v-model="amount"
        show-buttons
        :step="10"
        class="w-1/2"
      />
    </div>
    <p class="font-bold">
      {{ t('resources.move.price') }}: {{ price }}
    </p>
    <p class="font-bold">
      {{ t('resources.move.time') }}: {{ distance }}
    </p>
    <Button
      :label="t('common.move')"
      @click="move"
      class="w-1/2 self-center mt-4"
    />
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useResources } from '@/composables/useResources'
import type { DataMessage } from '@/types'
import type { Resource, ResourceMovePayload } from '@/types/Resources/index.interface'

interface Props {
  resource: Resource;
}
const props = defineProps<Props>()
const emits = defineEmits<{(e: 'close'): void}>()

const x = ref<number>(0)
const y = ref<number>(0)
const amount = ref<number>(0)
const message = ref<DataMessage | null>(null)

const { moveResource } = useResources()
const { t } = useI18n()

const distance = computed(() => {
  return ((props.resource.x-x.value)**2 + (props.resource.y-y.value)**2)**(0.5)
})
const price = computed(() => {
  return (props.resource.weight + props.resource.volume) * distance.value * amount.value / 1000
})

const move = () => {
  const payload: ResourceMovePayload = {
    toX: x.value,
    toY: y.value,
    amount: amount.value,
    resourceTypeId: props.resource.resourceTypeId,
    fromX: props.resource.x,
    fromY: props.resource.y
  }
  const { onFetchResponse, dataMessage } = moveResource(payload)
  onFetchResponse(() => {
    message.value = dataMessage.value

    setTimeout(() => {
      if (dataMessage.value?.status === 'success') {
        emits('close')
      }
    }, 2000)
  })
}
</script>

<style scoped>

</style>
