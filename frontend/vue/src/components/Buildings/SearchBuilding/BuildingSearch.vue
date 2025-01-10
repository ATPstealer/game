<template>
  <DataTable
    v-if="!isPending"
    class="p-4 w-full"
    paginator
    :pt="{paginator:
      {
        pagebutton: {class: 'paginator-button'},
        previouspagebutton: {class: 'paginator-button'},
        firstpagebutton: {class: 'paginator-button'},
        lastpagebutton: {class: 'paginator-button'},
        nextpagebutton: {class: 'paginator-button'},
      },
    }"
    :rows="10"
    :rows-per-page-options="[10, 25]"
    striped-rows
    :value="buildings as any[]"
  >
    <Column
      class="flex-2"
      field="title"
      :header="t('buildings.one')"
      :sortable="true"
    >
      <template #body="{data}">
        <span>{{ t(`buildings.types.${data?.buildingType.title.toLocaleLowerCase()}`) }}</span>
      </template>
    </Column>
    <Column
      class="flex-1"
      field="nickName"
      :header="t('common.owner')"
      :sortable="true"
    />
    <Column
      class="flex-1"
      field="square"
      :header="t('common.square')"
      :sortable="true"
    />
    <Column
      class="flex-1"
      field="x"
      header="X"
    />
    <Column
      class="flex-1"
      field="y"
      header="Y"
    />
  </DataTable>
</template>

<script setup lang="ts">
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import { onMounted, ref, toRefs, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { type BuildingWithData, type FindBuildingParams, usePostBuildingGet } from '@/gen'

interface Props {
  searchParams: FindBuildingParams;
}

const props = defineProps<Props>()
const { searchParams } = toRefs(props)

const buildings = ref<BuildingWithData[]>([])

const { t } = useI18n()

const { mutate, isPending } = usePostBuildingGet({
  mutation: {
    onSuccess: data => {
      buildings.value = data.data || []
    }
  }
})

watch(searchParams.value, () => {
  mutate({ data: { ...searchParams.value } })
})

onMounted(() => {
  mutate({ data: { ...searchParams.value } })
})

</script>

<style scoped>

</style>