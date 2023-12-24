<template>
  <div class="flex flex-col gap-4">
    <MessageBlock
      v-if="dataMessage"
      :message="dataMessage"
    />
    <p class="text-xl">
      Move
      <span class="font-bold">
        {{ resource.name }}
      </span>
      from
      <span class="font-bold">
        {{ resource.x }}:{{ resource.y }}
      </span>
      to
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
      <label class="font-bold text-xl">Amount:</label>
      <InputNumber
        v-model="amount"
        show-buttons
        :step="10"
        class="w-1/2"
      />
    </div>
    <p class="font-bold">
      Estimate price: {{ price }}
    </p>
    <p class="font-bold">
      Estimate time: {{ distance }}
    </p>
    <Button
      label="Move"
      @click="move"
      class="w-1/2 self-center mt-4"
    />
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import { computed, ref } from 'vue'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useResources } from '@/composables/useResources'
import type { Message } from '@/types'
import type { Resource } from '@/types/Resources/index.interface'

interface Props {
  resource: Resource;
}
const props = defineProps<Props>()
const emits = defineEmits<{(e: 'close'): void}>()

const x = ref<number>(0)
const y = ref<number>(0)
const amount = ref<number>(0)
const dataMessage = ref<Message>()
const { moveResource } = useResources()

const distance = computed(() => {
  return ((props.resource.x-x.value)**2 + (props.resource.y-y.value)**2)**(0.5)
})
const price = computed(() => {
  return (props.resource.weight + props.resource.volume) * distance.value * amount.value / 1000
})

const move = () => {
  const payload = {
    toX: x.value,
    toY: y.value,
    amount: amount.value,
    resourceTypeId: props.resource.resourceTypeId,
    fromX: props.resource.x,
    fromY: props.resource.y
  }
  const { onFetchResponse, message } = moveResource(payload)
  onFetchResponse(() => {
    dataMessage.value = message.value

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
