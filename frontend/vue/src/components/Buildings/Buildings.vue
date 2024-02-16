<template>
  <DataTable
    v-if="!isFetching"
    :value="buildings"
    size="small"
  >
    <Column
      :header="t(`buildings.one`)"
    >
      <template #body="{data}: {data: Building}">
        <p class="col-span-1 font-bold">
          <router-link :to="{name: `Building${data.buildingType.buildingGroup}Id`, params: {id: data._id}}" class="link">
            {{ t(`buildings.types.${data.buildingType.title.toLowerCase()}`) }} {{ data.level }}x{{ data.square }}
          </router-link>
        </p>
      </template>
    </Column>
    <Column
      :header="t(`map.cell`)"
    >
      <template #body="{data}: {data: Building}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
    <Column
      :header="t(`buildings.status`)"
    >
      <template #body="{data}: {data: Building}">
        {{ t(`status.${data.status.toLowerCase()}`) }}
      </template>
    </Column>
    <Column
      :header="t(`buildings.finish`)"
    >
      <template #body="{data}: {data: Building}">
        <span v-if="getTimeDiff(data.workEnd) > 0">
          {{ formatDuration(getTimeDiff(data.workEnd)) }}
        </span>
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
import { Building } from '@/types/Buildings/index.interface'
import { formatDuration } from '@/utils/formatDuration'
import { getTimeDiff } from '@/utils/getTimeDiff'

const { data: buildings, isFetching } = useGetData<Building[]>('/building/my')

const { t } = useI18n()
</script>
