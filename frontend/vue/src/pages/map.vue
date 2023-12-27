<template>
  <div v-if="!isFetching" class="flex flex-col">
    <div
      v-for="(row, rowIndex) in yArray"
      :key="rowIndex"
      class="flex justify-center"
    >
      <img
        v-for="(column, columnIndex) in xArray"
        :key="columnIndex"
        :src="getSrc(row, column)"
        class="border border-blue-950 border-solid"
        alt="Surface"
        @click="setCell(row, column)"
      />
    </div>
  </div>
  <Loading v-else />
  <Dialog
    v-model:visible="showModal"
    modal
    :header="`Cell ${cell?.x} x ${cell?.y}`"
  >
    <MapCell :square="cell?.square" :cell="cell" />
  </Dialog>
</template>

<script setup lang="ts">
import Dialog from 'primevue/dialog'
import { ref } from 'vue'
import Loading from '@/components/Common/Loading.vue'
import MapCell from '@/components/Map/MapCell.vue'
import { useMap } from '@/composables/useMap'

const cell = ref<any>()
const showModal = ref<boolean>(false)
const { getMap } = useMap()

const { data, onFetchResponse, xArray, yArray, isFetching } = getMap()

const getSrc = (row, column) => {
  const itemsInRow = data.value.filter(item => item.y === row)
  const oneCell = itemsInRow.filter(item => item.x === column)[0]

  return import.meta.env.VITE_MINIO_URL + oneCell.surfaceImagePath
}

const setCell = (row, column) => {
  const itemsInRow = data.value.filter(item => item.y === row)
  cell.value = itemsInRow.filter(item => item.x === column)[0]
  showModal.value = true
}

</script>

<style scoped>

</style>