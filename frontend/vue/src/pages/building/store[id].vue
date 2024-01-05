<template>
  <Layout>
    <template #help>
      {{ t(`buildings.store.help`) }}
    </template>
  <div v-if="!isFetching">
    <StoreBuilding :building="building" />
  </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { useGetData } from '@/composables/useGetData'
import type { Building } from '@/types/Buildings/index.interface'
import StoreBuilding from "@/components/Buildings/StoreBuilding/StoreBuilding.vue";
import Layout from "@/components/Common/Layout.vue";
import {useI18n} from "vue-i18n";

const route = useRoute()
const building = ref<Building>({} as Building)

const { data: myBuildings, onFetchResponse, isFetching } = useGetData<Building[]>('/building/my')
onFetchResponse(() => {
  building.value = myBuildings.value.find(item => item.id === Number(route.params.id)) as Building
})

const {t} = useI18n()
</script>

<style scoped>

</style>
