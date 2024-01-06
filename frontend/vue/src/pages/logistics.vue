<template>
  <Layout>
  <div v-if="!isFetching">
    <h2 class="font-bold text-2xl">
      Logistics
    </h2>
    <div class="grid grid-cols-5 font-bold">
      <h3 class="col-span-1">
        Resource
      </h3>
      <h3 class="col-span-1">
        From
      </h3>
      <h3 class="col-span-1">
        To
      </h3>
      <h3 class="col-span-1">
        Amount
      </h3>
      <h3 class="col-span-1">
        Finish
      </h3>
    </div>
    <div
      v-for="logistic in logistics"
      class="grid grid-cols-5 "
      :key="logistic.id"
    >
      <p class="col-span-1">
        {{ logistic.resourceName }}
      </p>
      <p class="col-span-1">
        {{ logistic.fromX }}x{{ logistic.fromY }}
      </p>
      <p class="col-span-1">
        {{ logistic.toX }}x{{ logistic.toY }}
      </p>
      <p class="col-span-1">
        {{ logistic.amount }}
      </p>
      <p class="col-span-1" v-if="getTimeDiff(logistic.workEnd)">
        {{ formatDuration(getTimeDiff(logistic.workEnd)) }}
      </p>
    </div>
  </div>
  <Loading v-else />
  </Layout>
</template>

<script setup lang="ts">
import Loading from '@/components/Common/Loading.vue'
import { useGetData } from '@/composables/useGetData'
import type { Logistic } from '@/types'
import { formatDuration } from '@/utils/formatDuration'
import { getTimeDiff } from '@/utils/getTimeDiff'
import Layout from "@/components/Common/Layout.vue";

const { data: logistics, isFetching } = useGetData<Logistic[]>('/resource/my_logistics')
</script>

<style scoped>

</style>
