<template>
  <Layout>
    <template #options>
      <Hiring v-if="!isFetching" :building="building" />
      <Button
        :label="t(`buildings.destroy`)"
        @click="confirmDelete($event, building._id)"
        class="w-max"
        size="small"
        severity="danger"
      />
    </template>
    <template #help>
      <p>{{ t(`buildings.hiring.help`) }}</p>
    </template>
    <div v-if="!isFetching">
      <ProductionBuilding :building="building" />
    </div>
  </Layout>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import { useConfirm } from 'primevue/useconfirm'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import Hiring from '@/components/Buildings/Hiring.vue'
import ProductionBuilding from '@/components/Buildings/ProductionBuilding/ProductionBuilding.vue'
import Layout from '@/components/Common/Layout.vue'
import { useBuildings } from '@/composables/useBuildings'
import { useGetData } from '@/composables/useGetData'
import type { Building } from '@/types/Buildings/index.interface'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const building = ref<Building>({} as Building)

const { destroyBuilding } = useBuildings()
const confirm = useConfirm()

const { data, onFetchResponse, isFetching } = useGetData<Building[]>(`/building/my?_id=${route.params.id}`)

onFetchResponse(() => {
  building.value = data.value[0]
})

const confirmDelete = (event: any, id: string) => {
  confirm.require({
    target: event.currentTarget,
    message: t('buildings.destroyConfirm'),
    icon: 'pi pi-info-circle',
    acceptClass: 'p-button-danger p-button-sm',
    acceptLabel: 'Да',
    rejectLabel: 'Нет',
    accept: () => destroy(id)
  })
}

const destroy = (id: string) => {
  const { onFetchResponse } = destroyBuilding(id)
  onFetchResponse(() => {
    router.push({ name: 'Buildings' })
  })
}

</script>

<style scoped>

</style>
