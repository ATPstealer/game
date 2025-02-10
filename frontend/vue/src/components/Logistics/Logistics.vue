<template>
  <h2 class="font-bold text-2xl">
    {{ t(`logistics.title`) }}
  </h2>
  <DataTable
    v-if="!isFetching"
    size="small"
    striped-rows
    :value="logistics"
  >
    <Column
      :header="t(`logistics.columns.resource`)"
    >
      <template #body="{data}: {data: LogisticWithData}">
        {{ t(`resources.types.${data?.resourceType?.name?.toLowerCase()}`) }}
      </template>
    </Column>
    <Column
      :header="t(`logistics.columns.from`)"
    >
      <template #body="{data}: {data: LogisticWithData}">
        {{ data.fromX }}x{{ data.fromY }}
      </template>
    </Column>
    <Column
      :header="t(`logistics.columns.to`)"
    >
      <template #body="{data}: {data: LogisticWithData}">
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
      <template #body="{data}: {data: LogisticWithData}">
        {{ getTimeDiff(data?.workEnd) > 0 ? formatDuration(getTimeDiff(data?.workEnd)) : '' }}
      </template>
    </Column>
  </DataTable>
  <Loading v-else />
</template>

<script setup lang="ts">
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import { computed, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import Loading from '@/components/Common/Loading.vue'
import { type LogisticWithData, useGetResourceMyLogistics } from '@/gen'
import { formatDuration } from '@/utils/formatDuration'
import { getTimeDiff } from '@/utils/getTimeDiff'

const { data: logisticsQuery, suspense, isFetching } = useGetResourceMyLogistics()
await suspense()
const logistics = computed(() => unref(logisticsQuery)?.data)

const { t } = useI18n()
</script>
