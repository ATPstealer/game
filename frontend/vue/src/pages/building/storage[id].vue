<template>
  <Layout>
    <template #options>
      <Hiring v-if="!isFetching" :building="building" />
    </template>
    <template #help>
      {{ t(`buildings.hiring.help`) }}
    </template>
    Storage
  </Layout>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import Hiring from '@/components/Buildings/Hiring.vue'
import Layout from '@/components/Common/Layout.vue'
import { useGetData } from '@/composables/useGetData'
import type { Building } from '@/types/Buildings/index.interface'

const route = useRoute()
const { t } = useI18n()
const building = ref<Building>({} as Building)

const { data: myBuildings, onFetchResponse, isFetching } = useGetData<Building[]>('/building/my')
onFetchResponse(() => {
  building.value = myBuildings.value.find(item => item._id === route.params.id) as Building
})
</script>
