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
      <template #body="{data}: {data: LandLord}">
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
import { computed, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import Loading from '@/components/Common/Loading.vue'
import { type LandLord, useGetMapMy } from '@/gen'

const { t } = useI18n()

const { data: landsQuery, isFetching, suspense } = useGetMapMy()
await suspense()
const lands = computed(() => unref(landsQuery)?.data)

</script>
