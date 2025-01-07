<template>
  <h2 class="font-bold text-2xl">
    {{ t(`lands.title`) }}
  </h2>
  <DataTable
    v-if="!isFetching"
    size="small"
    striped-rows
    :value="lands"
  >
    <Column
      :header="t(`map.cell`)"
    >
      <template #body="{data}: {data: Land}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
    <Column
      field="square"
      :header="t(`map.square`)"
    />
  </DataTable>
  <Loading v-else />
</template>

<script setup lang="ts">
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import { useI18n } from 'vue-i18n'
import Loading from '@/components/Common/Loading.vue'
import { useGetData } from '@/composables/useGetData'
import { Land } from '@/types'

const { data: lands, isFetching } = useGetData<Land[]>('/map/my')
const { t } = useI18n()
</script>
