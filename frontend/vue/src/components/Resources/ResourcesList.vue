<template>
  <DataTable
    :value="resources?.filter(resource => resource?.resourceType?.name)"
    size="small"
    striped-rows
  >
    <Column
      :header="t(`resources.columns.name`)"
    >
      <template #body="{data}: {data: Resource}">
        <span class="clickable-item" @click="openMoveResource(data)">
          {{ t(`resources.types.${data?.resourceType?.name?.toLowerCase()}`) }}
        </span>
      </template>
    </Column>
    <Column
      :header="t(`resources.columns.amount`)"
    >
      <template #body="{data}: {data: Resource}">
        <span @click="openSellResource(data)" class="clickable-item">
          {{ data.amount }}
        </span>
      </template>
    </Column>
    <Column
      :header="t(`map.cell`)"
    >
      <template #body="{data}: {data: Resource}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
  </DataTable>

  <Dialog
    v-model:visible="moveResourcesModal"
    modal
    :header="t('resources.move.header')"
    :style="{ width: '25rem' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :dismissable-mask="true"
  >
    <MoveResource
      :resource="currentResource"
      @close="onCloseMoveModal"
    />
  </Dialog>

  <Dialog
    v-model:visible="sellResourcesModal"
    modal
    :header="t('resources.sell.header')"
    :style="{ width: '25rem' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :dismissable-mask="true"
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
import { Resource } from '@/types/Resources/index.interface'

interface Props {
  resources: Resource[];
  execute: () => void;
}

const props = defineProps<Props>()

const { t } = useI18n()

const moveResourcesModal = ref<boolean>(false)
const sellResourcesModal = ref<boolean>(false)
const currentResource = ref<Resource>({} as Resource)

const openMoveResource = (resource: Resource) => {
  currentResource.value = resource
  moveResourcesModal.value = true
}

const openSellResource = (resource: Resource) => {
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
