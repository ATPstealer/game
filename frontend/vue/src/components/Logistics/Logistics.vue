<template>
  <h2 class="font-bold text-2xl">
    {{ t(`logistics.title`) }}
  </h2>
  <DataTable
    v-if="!isFetching"
    :value="logistics"
    size="small"
    striped-rows
  >
    <Column
      field="resourceName"
      :header="t(`logistics.columns.resource`)"
    />
    <Column
      :header="t(`logistics.columns.from`)"
    >
      <template #body="{data}: {data: Logistic}">
        {{ data.fromX }}x{{ data.fromY }}
      </template>
    </Column>
    <Column
      :header="t(`logistics.columns.to`)"
    >
      <template #body="{data}: {data: Logistic}">
        {{ data.toX }}x{{ data.toY }}
      </template>
    </Column>
    <Column
      field="amount"
      :header="t(`logistics.columns.amount`)"
    />
    <Column
      :header="t(`logistics.columns.finish`)"
    >
      <template #body="{data}: {data: Logistic}">
        {{ getTimeDiff(data.workEnd) > 0 ? formatDuration(getTimeDiff(data.workEnd)) : '' }}
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
import { Logistic } from '@/types'
import { formatDuration } from '@/utils/formatDuration'
import { getTimeDiff } from '@/utils/getTimeDiff'

const { data: logistics, isFetching } = useGetData<Logistic[]>('/resource/my_logistics')
const { t } = useI18n()
</script>
