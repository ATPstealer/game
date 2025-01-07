<template>
  <h2 class="mb-5 text-center">
    {{ t(`storages.title`) }}
  </h2>
  <DataTable
    v-if="!isFetching"
    size="small"
    striped-rows
    :value="storages"
  >
    <Column
      class="w-1/2"
      :header="t(`map.cell`)"
    >
      <template #body="{data}: {data: Storage}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
    <Column
      class="w-1/2"
      :header="t(`storages.columns.volume`)"
    >
      <template #body="{data}: {data: Storage}">
        <p :class="{'font-bold text-red-600': data.volumeOccupied > data.volumeMax}">
          {{ Math.trunc(data.volumeOccupied) }}/{{ data.volumeMax }}
        </p>
      </template>
    </Column>
  </DataTable>
  <Loading v-else />
</template>
<script setup lang="ts">
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import { useI18n } from 'vue-i18n'
import Loading from '@/components/Common/Loading.vue'
import { useGetData } from '@/composables/useGetData'
import { Storage } from '@/types'

const { data: storages, isFetching } = useGetData<Storage[]>('/storage/my')
const { t } = useI18n()
</script>
