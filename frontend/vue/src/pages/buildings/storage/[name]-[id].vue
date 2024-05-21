<template>
  <BuildingTemplate :building="building" :loading="isFetching">
    <template #building>
      Storage
    </template>
    <template #buildingHelp>
      <p>{{ t(`buildings.hiring.help`) }}</p>
    </template>
  </BuildingTemplate>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import BuildingTemplate from '@/components/Buildings/BuildingTemplate.vue'
import { useGetData } from '@/composables/useGetData'
import type { Building } from '@/types/Buildings/index.interface'

const route = useRoute()
const { t } = useI18n()

const building = ref<Building>({} as Building)

const { data, onFetchResponse, isFetching } = useGetData<Building[]>(`/building/my?_id=${route.params.id}`)

onFetchResponse(() => {
  building.value = data.value[0]
})
</script>
