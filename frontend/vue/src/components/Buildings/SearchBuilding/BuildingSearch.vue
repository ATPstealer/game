<template>
  <DataTable
    v-if="!isFetching"
    :value="buildings as any[]"
    striped-rows
    paginator
    :rows="10"
    :rows-per-page-options="[10, 25]"
    class="p-4 w-full"
    :pt="{paginator:
      {
        pagebutton: {class: 'paginator-button'},
        previouspagebutton: {class: 'paginator-button'},
        firstpagebutton: {class: 'paginator-button'},
        lastpagebutton: {class: 'paginator-button'},
        nextpagebutton: {class: 'paginator-button'},
      },
    }"
  >
    <Column
      field="title"
      :header="t('buildings.one')"
      :sortable="true"
      class="flex-2"
    >
      <template #body="{data}">
        <span>{{ t(`buildings.types.${data?.title.toLocaleLowerCase()}`) }}</span>
      </template>
    </Column>
    <Column
      field="nickName"
      :header="t('common.owner')"
      :sortable="true"
      class="flex-1"
    />
    <Column
      field="square"
      :header="t('common.square')"
      :sortable="true"
      class="flex-1"
    />
    <Column
      field="x"
      header="X"
      class="flex-1"
    />
    <Column
      field="y"
      header="Y"
      class="flex-1"
    />
  </DataTable>
</template>

<script setup lang="ts">
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import { toRefs } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBuildings } from '@/composables/useBuildings'
import type { SearchBuildingParams } from '@/types/Buildings/index.interface'

interface Props {
  searchParams: SearchBuildingParams;
}

const props = defineProps<Props>()
const { searchParams } = toRefs(props)

const { t } = useI18n()
const { getBuildings } = useBuildings()
const { data: buildings, isFetching } = getBuildings(searchParams)

</script>

<style scoped>

</style>