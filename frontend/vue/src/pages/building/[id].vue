<template>
  <div v-if="!isFetching">
    <ProductionBuilding :building="building" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import ProductionBuilding from '@/components/Buildings/ProductionBuilding/ProductionBuilding.vue'
import { useGetData } from '@/composables/useGetData'
import type { Building } from '@/types/Buildings/index.interface'

const route = useRoute()
const building = ref<Building>({} as Building)

const { data: myBuildings, onFetchResponse, isFetching } = useGetData<Building[]>('/building/my')
onFetchResponse(() => {
  building.value = myBuildings.value.find(item => item.id === Number(route.params.id)) as Building
})

</script>

<style scoped>

</style>