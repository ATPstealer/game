<template>
  <Layout>
    <template #options>
      <div class="slider">
        <span>X: {{ x }}</span>
        <div class="flex items-center gap-4">
          <span class="text-sm">{{ settings?.mapMinX }}</span>
          <Slider
            v-model="x"
            class="w-[150px]"
            :max="settings?.mapMaxX"
            :min="settings?.mapMinX"
            @change="value => x = value"
            @slideend="event => setParams({key: 'x', value: event.value})"
          />
          <span class="text-sm">{{ settings?.mapMaxX }}</span>
        </div>
      </div>
      <div class="slider">
        <span>Y: {{ y }}</span>
        <div class="flex items-center gap-4">
          <span class="text-sm">{{ settings?.mapMinY }}</span>
          <Slider
            v-model="y"
            class="w-[150px]"
            :max="settings?.mapMaxY"
            :min="settings?.mapMinY"
            @change="value => y = value"
            @slideend="event => setParams({key: 'y', value: event.value})"
          />
          <span class="text-sm">{{ settings?.mapMaxY }}</span>
        </div>
      </div>
      <Dropdown
        v-if="resourcesTypes?.length"
        v-model="currentResource"
        :option-label="event => t(`resources.types.${event.name.toLowerCase()}`)"
        option-value="id"
        :options="resourcesTypes"
        @change="event => setParams({key: 'resourceTypeId', value: event.value})"
      />
      <Button :label="t('common.reset')" @click="reset" />
      <Button
        class="mb-8"
        :label="t('orders.create.header')"
        severity="info"
        @click="sellResourcesModal= true"
      />
      <Dialog
        v-model:visible="sellResourcesModal"
        :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
        :dismissable-mask="true"
        :header="t('orders.create.header')"
        modal
        :style="{ width: '25rem' }"
      >
        <CreateOrderModal
          @close="sellResourcesModal = false; setParams({key: 'trigger', value: 0})"
        />
      </Dialog>
    </template>
    <MarketPage :search-params="params" />
  </Layout>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import Dropdown from 'primevue/dropdown'
import Slider from 'primevue/slider'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'
import CreateOrderModal from '@/components/Market/CreateOrderModal.vue'
import MarketPage from '@/components/Market/MarketPage.vue'
import { useGetData } from '@/composables/useGetData'
import { Coords } from '@/types'
import type { MarketParams, ResourceType } from '@/types/Resources/index.interface'

const x = ref<number>()
const y = ref<number>()
const currentResource = ref<number>(0)
const params = ref<MarketParams>({} as MarketParams)
const sellResourcesModal = ref<boolean>(false)
const resourcesTypes = ref<ResourceType[]>([{ id: 0, name: 'All' } as ResourceType])

const { data: settings } = useGetData<Record<Coords, number>>('/settings')
const { data: resourcesTypesData, onFetchResponse } = useGetData<ResourceType[]>('/resource/types')
const { t } = useI18n()

onFetchResponse(() => {
  resourcesTypes.value = [...resourcesTypes.value, ...resourcesTypesData.value]
})

const setParams = ({ key, value }: {key: string; value: number}) => {
  params.value[key] = value
  delete params.value.trigger

  if (!currentResource.value) {
    delete params.value.resourceTypeId
  }
}

const reset = () => {
  x.value = undefined
  delete params.value.x
  y.value = undefined
  delete params.value.y
}

</script>

<style scoped>

</style>