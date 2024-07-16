<template>
  <BuildingTemplate :building="building" :loading="isFetching">
    <template #building>
      <StorageBuilding :building="building" />
    </template>
    <template #buildingHelp>
      <p>{{ t(`buildings.hiring.help`) }}</p>
    </template>
  </BuildingTemplate>
</template>
<script setup lang="ts">
import { provide, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import BuildingTemplate from '@/components/Buildings/BuildingTemplate.vue'
import StorageBuilding from '@/components/Buildings/StorageBuilding/StorageBuilding.vue'
import { useGetData } from '@/composables/useGetData'
import type { Building } from '@/types/Buildings/index.interface'

const route = useRoute()
const { t } = useI18n()

const building = ref<Building>({} as Building)

const { data, onFetchResponse, isFetching, execute } = useGetData<Building[]>(`/building/my?_id=${route.params.id}`)

onFetchResponse(() => {
  building.value = data.value[0]
})

provide('execute', execute)
</script>
