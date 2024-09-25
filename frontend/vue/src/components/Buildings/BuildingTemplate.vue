<template>
  <Layout>
    <template #options>
      <div class="flex flex-col gap-8">
        <Hiring v-if="!loading" :building="building" />
        <Equipment v-if="!loading" :building="building" />
        <Divider />
        <div class="flex flex-col gap-4">
          <Button
            class="w-max"
            :label="t(`buildings.destroy`)"
            severity="danger"
            size="small"
            @click="confirmDelete($event, building._id)"
          />
          <Button
            class="w-max"
            :label="t(`buildings.stop work`)"
            severity="danger"
            size="small"
            @click="confirmStopWork($event, building._id)"
          />
        </div>
      </div>
    </template>
    <template #help>
      <slot name="buildingHelp" />
    </template>
    <div v-if="!loading" class="p-4">
      <slot name="building" />
    </div>
    <Loading v-else />
  </Layout>
</template>

<script setup lang="ts">

import Button from 'primevue/button'
import Divider from 'primevue/divider'
import { useConfirm } from 'primevue/useconfirm'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import Equipment from '@/components/Buildings/Equipment.vue'
import Hiring from '@/components/Buildings/Hiring.vue'
import Layout from '@/components/Common/Layout.vue'
import Loading from '@/components/Common/Loading.vue'
import { useBuildings } from '@/composables/useBuildings'
import type { Building } from '@/types/Buildings/index.interface'

interface Props {
  building: Building;
  loading: boolean;
}

defineProps<Props>()

const router = useRouter()
const { t } = useI18n()
const { destroyBuilding, stopWork } = useBuildings()
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

const confirmStopWork = (event: any, id: string) => {
  confirm.require({
    target: event.currentTarget,
    message: t('common.confirm'),
    icon: 'pi pi-info-circle',
    acceptClass: 'p-button-danger p-button-sm',
    acceptLabel: t('common.yes'),
    rejectLabel: t('common.no'),
    accept: () => stop(id)
  })
}

const stop = (id: string) => {
  const payload = {
    buildingId: id
  }
  const { onFetchResponse } = stopWork(payload)
  onFetchResponse(() => {})
}

</script>

<style scoped>

</style>