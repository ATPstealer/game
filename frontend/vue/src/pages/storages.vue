<template>
  <Layout>
    <div v-if="!isFetching">
      <h2 class="font-bold text-2xl">
        Storages
      </h2>
      <div class="grid grid-cols-3 gap-2 font-bold">
        <h3 class="col-span-1">
          Cell
        </h3>
        <h3 class="col-span-1">
          Volumes
        </h3>
      </div>

      <div
        v-for="storage in storages"
        class="grid grid-cols-3 "
        :key="storage.id"
      >
        <p class="col-span-1">
          {{ storage.x }}:{{ storage.y }}
        </p>
        <p :class="{'font-bold text-red-600': storage.volumeOccupied > storage.volumeMax}">
          {{ storage.volumeOccupied }}/{{ storage.volumeMax }}
        </p>
      </div>
      <p class="text-white" />
    </div>
    <Loading v-else />
  </Layout>
</template>

<script setup lang="ts">
import Layout from '@/components/Common/Layout.vue'
import Loading from '@/components/Common/Loading.vue'
import { useGetData } from '@/composables/useGetData'
import type { Storage } from '@/types'

const { data: storages, isFetching } = useGetData<Storage[]>('/storage/my')
</script>

<style scoped>

</style>
