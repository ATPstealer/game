<template>
  <BuildingTemplate :building="building" :loading="isFetching">
    <template #building />
    <template #buildingHelp>
      <p>{{ t(`buildings.store.help`) }}</p>
      <hr />
      <p>{{ t(`buildings.hiring.help`) }}</p>
    </template>
  </BuildingTemplate>
</template>

<script setup lang="ts">
import { computed, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import BuildingTemplate from '@/components/Buildings/BuildingTemplate.vue'
import { useGetBuildingMy } from '@/gen'

const route = useRoute()
const { t } = useI18n()

const { data: storageBuildingQuery, suspense, isFetching } = useGetBuildingMy({ _id: route.params.id })
await suspense()
const building = computed(() => unref(storageBuildingQuery)?.data?.find(item => item._id === route.params.id))

</script>

<style scoped>

</style>
