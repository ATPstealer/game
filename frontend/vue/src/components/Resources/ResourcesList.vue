<template>
  <DataTable
    size="small"
    striped-rows
    :value="resources?.filter(resource => resource?.resourceType?.name)"
  >
    <Column
      :header="t(`resources.columns.name`)"
    >
      <template #body="{data}: {data: ResourceWithData}">
        <span class="clickable-item" @click="openMoveResource(data)">
          {{ t(`resources.types.${data?.resourceType?.name?.toLowerCase()}`) }}
        </span>
      </template>
    </Column>
    <Column
      :header="t(`resources.columns.amount`)"
    >
      <template #body="{data}: {data: ResourceWithData}">
        <span class="clickable-item" @click="openSellResource(data)">
          {{ data.amount }}
        </span>
      </template>
    </Column>
    <Column
      :header="t(`map.cell`)"
    >
      <template #body="{data}: {data: ResourceWithData}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
  </DataTable>

  <Dialog
    v-model:visible="moveResourcesModal"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :dismissable-mask="true"
    :header="t('resources.move.header')"
    modal
    :style="{ width: '25rem' }"
  >
    <MoveResource
      :resource="currentResource"
      @close="onCloseMoveModal"
    />
  </Dialog>

  <Dialog
    v-model:visible="sellResourcesModal"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :dismissable-mask="true"
    :header="t('resources.sell.header')"
    modal
    :style="{ width: '25rem' }"
  >
    <CreateOrderModal
      :resource="currentResource"
      @close="onCloseOrderModal"
    />
  </Dialog>
</template>
<script setup lang="ts">
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import CreateOrderModal from '@/components/Market/CreateOrderModal.vue'
import MoveResource from '@/components/Resources/MoveResource.vue'
import type { ResourceWithData } from '@/gen'

interface Props {
  resources: ResourceWithData[] | undefined;
  execute: () => void;
}

const props = defineProps<Props>()

const { t } = useI18n()

const moveResourcesModal = ref<boolean>(false)
const sellResourcesModal = ref<boolean>(false)
const currentResource = ref<ResourceWithData>({} as ResourceWithData)

const openMoveResource = (resource: ResourceWithData) => {
  currentResource.value = resource
  moveResourcesModal.value = true
}

const openSellResource = (resource: ResourceWithData) => {
  currentResource.value = resource
  sellResourcesModal.value = true
}

const onCloseOrderModal = () => {
  sellResourcesModal.value = false
  props.execute()
}

const onCloseMoveModal = () => {
  moveResourcesModal.value = false
  props.execute()
}
</script>
