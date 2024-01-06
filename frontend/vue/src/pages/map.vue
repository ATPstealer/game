<template>
  <Layout>
    <div v-if="!isMapFetching && !isSettingsFetching" class="flex flex-col">
      <div
        v-for="(row, rowIndex) in yArray"
        :key="`row-${rowIndex}`"
        class="flex justify-center"
      >
        <img
          v-for="(column, columnIndex) in xArray"
          :key="`column-${columnIndex}`"
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
  </Layout>
</template>

<script setup lang="ts">
import Dialog from 'primevue/dialog'
import { ref } from 'vue'
import Layout from '@/components/Common/Layout.vue'
import Loading from '@/components/Common/Loading.vue'
import MapCell from '@/components/Map/MapCell.vue'
import { useGetData } from '@/composables/useGetData'
import { useMap } from '@/composables/useMap'
import type { Coords } from '@/types'

const cell = ref<any>()
const showModal = ref<boolean>(false)
const xArray = ref<number[]>([])
const yArray = ref<number[]>([])

const { getMap } = useMap()

const { data: settings, onFetchResponse, isFetching: isSettingsFetching } = useGetData<Record<Coords, number>>('/settings')
const { data: map, isFetching: isMapFetching } = getMap()

onFetchResponse(() => {
  for (let i = settings.value.mapMinX; i <= settings.value.mapMaxX; i++) {
    xArray.value.push(i)
  }
  for (let i = settings.value.mapMinY; i <= settings.value.mapMaxY; i++) {
    yArray.value.push(i)
  }
})

const getSrc = (row, column) => {
  const itemsInRow = map.value.filter(item => item.y === row)
  const oneCell = itemsInRow.filter(item => item.x === column)[0]

  return import.meta.env.VITE_MINIO_URL + oneCell.surfaceImagePath
}

const setCell = (row, column) => {
  const itemsInRow = map.value.filter(item => item.y === row)
  cell.value = itemsInRow.filter(item => item.x === column)[0]
  showModal.value = true
}

</script>

<style scoped>

</style>
