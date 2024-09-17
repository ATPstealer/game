<template>
  <DataTable
    v-if="!isFetching"
    size="small"
    :value="buildings"
  >
    <Column
      :header="t(`buildings.one`)"
    >
      <template #body="{data}: {data: Building}">
        <p class="col-span-1 font-bold">
          <router-link class="link" :to="{name: `Buildings${data.buildingType.buildingGroup}NameId`, params: {id: data._id, name: data.buildingType.title.toLowerCase()}}">
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
        {{ t(`status.${data?.status?.toLowerCase()}`) }}
      </template>
    </Column>
    <Column
      :header="t(`buildings.finish`)"
    >
      <template #body="{data}: {data: Building}">
        <div class="min-w-[96px] max-w-[96px]">
          {{ getTime(data.workEnd, 'Достроено') }}
        </div>
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
import { useTimer } from '@/composables/useTimer'
import type { Building } from '@/types/Buildings/index.interface'

const { data: buildings, isFetching } = useGetData<Building[]>('/building/my')
const { getTime } = useTimer()

const { t } = useI18n()
</script>
