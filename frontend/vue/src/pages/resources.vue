<template>
  <div v-if="!isFetching">
    <h2 class="font-bold text-2xl">
      Resources
    </h2>
    <div class="grid grid-cols-4 font-bold">
      <h3 class="col-span-1">
        Resource
      </h3>
      <h3 class="col-span-1">
        Amount
      </h3>
      <h3 class="col-span-1">
        Cell
      </h3>
      <h3 class="col-span-1">
        Sell
      </h3>
    </div>
    <div
      v-for="resource in resources"
      class="grid grid-cols-4"
      :key="resource.id"
    >
      <p class="col-span-1" @click="openMoveResource(resource)">
        {{ resource.name }}
      </p>
      <p class="col-span-1">
        {{ resource.amount }}
      </p>
      <p class="col-span-1">
        {{ resource.x }}x{{ resource.y }}
      </p>
      <p class="col-span-1" @click="openSellResource(resource)">
        sell
      </p>
    </div>
  </div>
  <Loading v-else />

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
      @close="moveResourcesModal = false"
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
      @close="sellResourcesModal = false"
    />
  </Dialog>
</template>

<script setup lang="ts">

import Dialog from 'primevue/dialog'
import { ref } from 'vue'
import Loading from '@/components/Common/Loading.vue'
import CreateOrderModal from '@/components/Market/CreateOrderModal.vue'
import MoveResource from '@/components/Resources/MoveResource.vue'
import { useGetData } from '@/composables/useGetData'
import type { Resource } from '@/types/Resources/index.interface'

const { data: resources, isFetching } = useGetData<Resource[]>('/resource/my')

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

</script>

<style scoped>

</style>
