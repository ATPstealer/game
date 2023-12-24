<template>
  <div v-if="!isFetching">
    <h2 class="font-bold text-2xl">
      Buildings
      <router-link
        class="text-blue-500 hover:text-blue-700"
        :to="{name: 'ConstructBuilding', query: {x: 0, y: 0}}"
      >
        Construct
      </router-link>
    </h2>
    <div class="text-black">
      <div class="grid grid-cols-4 gap-2 font-bold">
        <h3 class="col-span-1">
          Building
        </h3>
        <h3 class="col-span-1">
          Cell
        </h3>
        <h3 class="col-span-1">
          Status
        </h3>
        <h3 class="col-span-1">
          Finish
        </h3>
      </div>

      <div
        v-for="building in buildings"
        class="grid grid-cols-4 "
        :key="building.id"
      >
        <p class="col-span-1 font-bold text-blue-500 hover:text-blue-700">
          <router-link :to="{name: 'BuildingId', params: {id: building.id}}">
            {{ building.title }} {{ building.level }}x{{ building.square }}
          </router-link>
        </p>
        <p class="col-span-1">
          {{ building.x }}x{{ building.y }}
        </p>
        <p class="col-span-1">
          {{ building.status }}
        </p>
        <p class="col-span-1 ml-2" v-if="getTimeDiff(building.workEnd) > 0">
          {{ formatDuration(getTimeDiff(building.workEnd)) }}
        </p>
      </div>
    </div>
  </div>
  <Loading v-else />
</template>

<script setup lang="ts">
import Loading from '@/components/Common/Loading.vue'
import { useGetData } from '@/composables/useGetData'
import type { Building } from '@/types/Buildings/index.interface'
import { formatDuration } from '@/utils/formatDuration'
import { getTimeDiff } from '@/utils/getTimeDiff'

const { data: buildings, isFetching } = useGetData<Building[]>('/building/my')
</script>

<style scoped>

</style>