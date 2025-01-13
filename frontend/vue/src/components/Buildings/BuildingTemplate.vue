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
            @click="confirmDelete($event, building!._id)"
          />
          <Button
            class="w-max"
            :label="t(`buildings.stop work`)"
            severity="danger"
            size="small"
            @click="confirmStopWork($event, building!._id)"
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
import {
  type BuildingWithData,
  type PostBuildingStopWorkMutationRequest,
  useDeleteBuildingDestroy,
  usePostBuildingStopWork
} from '@/gen'

interface Props {
  building: BuildingWithData | undefined;
  loading: boolean;
}

defineProps<Props>()

const router = useRouter()
const { t } = useI18n()
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

const destroyMutation = useDeleteBuildingDestroy({
  mutation: {
    onSuccess: () => {
      router.push({ name: 'Buildings' })
    }
  }
})

const destroy = (id: string) => {
  destroyMutation.mutate({ params: { _id: id } })
}

const stopMutation = usePostBuildingStopWork()

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
  // TODO: неверная типизация payload в сваггере
  stopMutation.mutate({ data: { ...payload } as PostBuildingStopWorkMutationRequest })
}

</script>

<style scoped>

</style>