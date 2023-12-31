<template>
  <DataTable
      v-if="!isFetching"
      :value="buildings"
      size="small"
  >
    <Column
        :header="t(`buildings.one`)"
    >
      <template #body="{data}">
        <p class="col-span-1 font-bold text-blue-500 hover:text-blue-700">
          <router-link :to="{name: `Building${data.buildingGroup}Id`, params: {id: data.id}}">
            {{ data.title }} {{ data.level }}x{{ data.square }}
          </router-link>
        </p>
      </template>
    </Column>
    <Column
        :header="t(`map.cell`)"
    >
      <template #body="{data}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
    <Column
        field="status"
        :header="t(`buildings.status`)"
    />
    <Column
        :header="t(`buildings.finish`)"
    >
      <template #body="{data}">
        <span v-if="getTimeDiff(data.workEnd) > 0">
        {{ formatDuration(getTimeDiff(data.workEnd)) }}
        </span>
      </template>
    </Column>
  </DataTable>
  <Loading v-else/>
</template>

<script setup lang="ts">
import {getTimeDiff} from "@/utils/getTimeDiff";
import {formatDuration} from "@/utils/formatDuration";
import Loading from "@/components/Common/Loading.vue";
import {useGetData} from "@/composables/useGetData";
import {Building} from "@/types/Buildings/index.interface";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import {useI18n} from "vue-i18n";

const {data: buildings, isFetching} = useGetData<Building[]>('/building/my')

const {t} = useI18n()
</script>
