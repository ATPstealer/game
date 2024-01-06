<template>
  <div v-if="!isFetching">
    <MessageBlock v-if="message" :message="message" />
    <p>
      Each free slot of land costs $10 more than the previous one.
    </p>
    <p>Thus, the first one costs $10, the hundredth one costs $1000</p>
    <div class="font-bold">
      Free space: {{ freeSquare }} Ares
    </div>
    <div class="font-bold">
      Occupied land: {{ landOccupied }} Ares
    </div>
    <div class="flex flex-col gap-4 w-1/2 my-4">
      <p>Buy square:</p>
      <InputNumber v-model="buySquare" show-buttons />
      <p v-if="price">
        Estimate price {{ price }}
      </p>
      <Button label="Buy" @click="buy" />
    </div>
    <RouterLink :to="`/construct-building?x=${cell.x}&y=${cell.y}`" class="font-bold text-blue-500 hover:text-blue-700">
      Construct building
    </RouterLink>

    <div class="font-bold">
      <p class="mt-4">
        Pollution: {{ cell.pollution }}
      </p>
      <p>
        Population: {{ cell.population }}
      </p>
      <p>
        Education: {{ cell.education }}
      </p>
      <p>
        Crime: {{ cell.crime }}
      </p>
      <p>
        Medicine: {{ cell.medicine }}
      </p>
      <p>
        ElementarySchool: {{ cell.elementarySchool }}
      </p>
      <p>
        HigherSchool: {{ cell.higherSchool }}
      </p>
    </div>
    <div v-if="cellOwners">
      <p>
        Landlords:
      </p>
      <p v-for="owner in cellOwners" :key="owner.nickName">
        {{ owner.nickName }}: {{ owner.square }} Are
      </p>
    </div>
  </div>
  <Loading v-else />
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import { computed, ref } from 'vue'
import Loading from '@/components/Common/Loading.vue'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useMap } from '@/composables/useMap'
import type { Message } from '@/types'
import type { Cell } from '@/types/Map/index.interface'

interface Props {
  cell: Cell;
  square: number;
}

const props = defineProps<Props>()

const buySquare = ref<number>(0)

const { getCellOwners, buyCellSquare } = useMap()
const { data: cellOwners, onFetchResponse, isFetching } = getCellOwners({ x: props.cell.x, y: props.cell.y })
const message = ref<Message | null>(null)

const freeSquare = computed(() => {
  let value = props.square
  if (cellOwners.value) {
    for (const owner of cellOwners.value) {
      value -= owner.square
    }
  }

  return value
})

const landOccupied = computed(() => {
  let value = 0
  if (cellOwners.value) {
    for (const owner of cellOwners.value) {
      value += owner.square
    }
  }

  return value
})

const price = computed(() => {
  return 10 * (landOccupied.value * 2 + 1 + buySquare.value ) * buySquare.value / 2
})

const buy = () => {
  if (!buySquare.value) {
    return
  }

  const { dataMessage, onFetchResponse } = buyCellSquare({ x: props.cell.x, y: props.cell.y, square: buySquare.value })
  onFetchResponse(() => {
    message.value = dataMessage.value
  })
}
</script>

<style scoped>

</style>