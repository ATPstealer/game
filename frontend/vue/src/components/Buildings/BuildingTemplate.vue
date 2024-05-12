<template>
  <Layout>
    <template #options>
      <Hiring v-if="!loading" :building="building" />
      <Button
        :label="t(`buildings.destroy`)"
        @click="confirmDelete($event, building._id)"
        class="w-max"
        size="small"
        severity="danger"
      />
    </template>
    <template #help>
      <slot name="buildingHelp" />
    </template>
    <div v-if="!loading">
      <slot name="building" />
    </div>
    <Loading v-else />
  </Layout>
</template>

<script setup lang="ts">

import Button from 'primevue/button'
import { useConfirm } from 'primevue/useconfirm'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import Hiring from '@/components/Buildings/Hiring.vue'
import Layout from '@/components/Common/Layout.vue'
import Loading from '@/components/Common/Loading.vue'
import { useBuildings } from '@/composables/useBuildings'
import type { Building } from '@/types/Buildings/index.interface'

interface Props {
  building: Building;
  loading: boolean;
}

const router = useRouter()
const { t } = useI18n()

defineProps<Props>()

const { destroyBuilding } = useBuildings()
const confirm = useConfirm()

const confirmDelete = (event: any, id: string) => {
  confirm.require({
    target: event.currentTarget,
    message: t('buildings.destroyConfirm'),
    icon: 'pi pi-info-circle',
    acceptClass: 'p-button-danger p-button-sm',
    acceptLabel: t('common.yes'),
    rejectLabel: t('common.no'),
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