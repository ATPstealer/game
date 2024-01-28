<template>
  <Layout>
    <template #options>
      <div class="flex flex-col gap-4">
        <div
          v-for="option in filters"
          :key="`filter-${option}`"
        >
          <RadioButton
            v-model="filter"
            :value="option"
            :input-id="`${option}-name`"
          />
          <label :for="`${option}-name`" class="ml-2 capitalize">{{ option }}</label>
        </div>
      </div>
    </template>
    <div v-if="!isMapFetching && !isSettingsFetching && !recalculate" class="flex flex-col gap-8 items-center justify-center">
      <div class="flex leading-[0]">
        <div
          v-for="(row, rowIndex) in yArray"
          :key="`row-${rowIndex}`"
          class="flex flex-col-reverse "
        >
          <div
            v-for="(column, columnIndex) in xArray"
            :key="`column-${columnIndex}`"
            class="border border-blue-950 border-solid"
            :style="{'background-color': `rgba(${color}, ${getMapFilter(row, column)})`}"
            @click="setCell(row, column)"
          >
            <img
              :src="getSrc(row, column)"
              alt="Surface"
              class="relative z-[-1]"
            />
          </div>
        </div>
      </div>
      <div class="legend flex items-center relative">
        <span class="absolute -bottom-6">{{ computedFilter?.min }}</span>
        <span class="absolute -bottom-6 left-1/2 -translate-x-1/2" v-if="computedFilter?.max-computedFilter?.min">
          {{ Math.floor((computedFilter?.max-computedFilter?.min)/2) }}
        </span>
        <span class="absolute -bottom-6 left-full -translate-x-full">{{ Math.floor(computedFilter?.max) }}</span>
        <div
          class="h-10 w-10"
          v-for="item in 10"
          :key="`legend-${item}`"
          :style="{'background-color': `rgba(${color}, ${item * 0.1 * 0.7})`}"
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
import RadioButton from 'primevue/radiobutton'
import { computed, ref, watch } from 'vue'
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
const filter = ref<string>('pollution')
const recalculate = ref<boolean>(false)
const color = ref<string>('175, 27, 27')

const filters = ['pollution', 'population', 'education', 'crime', 'medicine', 'averageSalary']

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
  const oneCell = map.value.filter(item => item.x === row && item.y === column)[0]

  return import.meta.env.VITE_MINIO_URL + oneCell.surfaceImagePath
}

const setCell = (row, column) => {
  cell.value = map.value.filter(item => item.x === row && item.y === column)[0]
  showModal.value = true
}

const computedFilter = computed(() => {
  if (!map.value?.length) {
    return
  }
  const propertyArray = map.value.map(item => item[filter.value])

  return {
    min: Math.min(...propertyArray),
    max: Math.max(...propertyArray)
  }
})

const getMapFilter = (row, column) => {
  const cell = map.value.filter(item => item.y === column && item.x === row)[0]
  const { min, max } = computedFilter.value
  const range = max - min

  return ((cell[filter.value] - min) / range) * 0.6
}

watch(filter, () => {
  color.value = '175, 27, 27'
  if (!['pollution', 'crime'].includes(filter.value)) {
    color.value = '30, 137, 19'
  }
  recalculate.value = true
  setTimeout(() => {
    recalculate.value = false
  }, 0)
})

</script>

<style scoped>

</style>
