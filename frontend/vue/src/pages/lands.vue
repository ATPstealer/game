<template>
  <div v-if="!isFetching">
    <h2 class="font-bold text-2xl">
      Lands
    </h2>
    <div class="grid grid-cols-2 gap-2">
      <h3 class="col-span-1">
        Cell
      </h3>
      <h3 class="col-span-1">
        Square
      </h3>
    </div>
    <div
      v-for="(land, index) in lands"
      class="grid grid-cols-2 "
      :key="index"
    >
      <p class="col-span-1">
        {{ land.x }}x{{ land.y }}
      </p>
      <p class="col-span-1">
        {{ land.square }} Are
      </p>
    </div>
  </div>
  <Loading v-else />
</template>

<script setup lang="ts">
import Loading from '@/components/Common/Loading.vue'
import { useGetData } from '@/composables/useGetData'
import type { Land } from '@/types'

const { data: lands, isFetching } = useGetData<Land[]>('/map/my')
</script>

<style scoped>

</style>