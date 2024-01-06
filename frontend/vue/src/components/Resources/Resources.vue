<template>
  <DataTable
      v-if="!isFetching"
      :value="resources"
      size="small"
      striped-rows
  >
    <Column
        :header="t(`resources.columns.name`)"
    >
      <template #body="{data}">
      <span class="col-span-1 font-bold text-blue-500 hover:text-blue-700" @click="openMoveResource(data)">
        {{ data.name }}
      </span>
      </template>
    </Column>
    <Column
        :header="t(`resources.columns.amount`)"
    >
      <template #body="{data}">
      <span @click="openSellResource(data)" class="col-span-1 font-bold text-blue-500 hover:text-blue-700">
        {{ data.amount }}
      </span>
      </template>
    </Column>
    <Column
        :header="t(`map.cell`)"
    >
      <template #body="{data}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
  </DataTable>
  <Loading v-else/>

  <Dialog
      v-model:visible="moveResourcesModal"
      modal
      header="Move resource"
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
      header="Sell resource"
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
import MoveResource from "@/components/Resources/MoveResource.vue";
import Dialog from "primevue/dialog";
import Loading from "@/components/Common/Loading.vue";
import CreateOrderModal from "@/components/Market/CreateOrderModal.vue";
import {useGetData} from "@/composables/useGetData";
import {Resource} from "@/types/Resources/index.interface";
import {ref} from "vue";
import {useI18n} from "vue-i18n";
import DataTable from "primevue/datatable";
import Column from "primevue/column";

const {data: resources, isFetching, execute} = useGetData<Resource[]>('/resource/my')
const {t} = useI18n()

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
  execute()
}

const onCloseMoveModal = () => {
  moveResourcesModal.value = false
  execute()
}
</script>
