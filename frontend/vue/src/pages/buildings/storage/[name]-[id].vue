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
import { computed, provide, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import BuildingTemplate from '@/components/Buildings/BuildingTemplate.vue'
import StorageBuilding from '@/components/Buildings/StorageBuilding/StorageBuilding.vue'
import { useGetBuildingMy } from '@/gen'

const route = useRoute()
const { t } = useI18n()

const { data: storageBuildingQuery, suspense, isFetching, refetch } = useGetBuildingMy({ _id: route.params.id })
await suspense()
const building = computed(() => unref(storageBuildingQuery)?.data?.find(item => item._id === route.params.id))

provide('execute', refetch)
</script>
