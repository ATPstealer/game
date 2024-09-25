<template>
  <Layout class="mt-8" :show-options-prop="false">
    <Card class="max-w-screen-xl mx-auto xl:mt-16">
      <template #content>
        <div v-if="!isFetching" class="flex flex-col gap-4">
          <MessageBlock
            v-if="messageData"
            class="mb-8"
            v-bind="messageData"
            @close-message="messageData = undefined"
          />
          <div
            class="flex flex-col gap-8 xl:flex-row xl:gap-4 xl:items-center"
            :class="{'xl:mt-28': !messageData}"
          >
            <FloatLabel>
              <Dropdown
                id="resource"
                v-model="resource"
                class="w-full xl:max-w-[200px] xl:min-w-[200px]"
                :option-label="e => t(`resources.types.${e.name?.toLowerCase()}`)"
                option-value="id"
                :options="computedResources"
                :placeholder="t('resources.choose')"
              />
              <label for="resource">{{ t('resources.choose') }}</label>
            </FloatLabel>
            <FloatLabel>
              <Dropdown
                id="from"
                v-model="from"
                class="w-full xl:max-w-[200px] xl:min-w-[200px]"
                :disabled="!resource || computedFrom.length === 1"
                :option-label="e => `${e.cell} (${e.amount})`"
                :option-value="e => e.cell"
                :options="computedFrom"
                :placeholder="t('common.from')"
              />
              <label for="from">{{ t('common.from') }}</label>
            </FloatLabel>
            <FloatLabel>
              <Dropdown
                id="to"
                v-model="to"
                class="w-full xl:max-w-[150px] xl:min-w-[150px]"
                :disabled="!from"
                filter
                :options="computedTo"
                :placeholder="t('common.to')"
              />
              <label for="to">{{ t('common.to') }}</label>
            </FloatLabel>
            <FloatLabel>
              <InputNumber
                id="amount"
                v-model="amount"
                class="w-full xl:w-auto"
                :disabled="!to"
                :max="computedFrom?.find(item => item.cell === from)?.amount"
                show-buttons
                @input="e => amount = +e!.value"
              />
              <label for="amount">{{ t('common.amount') }}</label>
            </FloatLabel>
            <div id="table" />
            <div
              v-tooltip.bottom="!hub ? `${t('logistics.chooseHub')}` : ''"
              class="w-max"
              tabindex="0"
            >
              <Button
                :disabled="!hub"
                :label="t('logistics.send')"
                severity="info"
                type="submit"
                @click="create"
              />
            </div>
          </div>
          <div class="flex flex-col gap-4">
            <p><span class="font-bold">{{ t('logistics.reqCapacity') }}</span>: <span v-if="capacity >= 0" class="font-bold">{{ capacity }}</span></p>
            <p><span class="font-bold">{{ t('logistics.cost') }}</span>: <span v-if="price" class="font-bold">{{ capacity * price }}$</span></p>
          </div>
          <!--          <Teleport :disabled="isDesktop" to="#table">-->
          <DataTable
            :row-class="rowClass"
            selection-mode="single"
            :value="hubs"
            @row-click="onRowClick"
          >
            <template #header>
              <h3>{{ t('logistics.hub') }}</h3>
            </template>
            <template #empty>
              <span v-if="capacity >= 0">{{ t('logistics.noHubs') }}</span>
              <span v-else>{{ t('common.completeSteps') }}</span>
            </template>
            <Column field="capacity" :header="t('logistics.capacity')" />
            <Column field="speed" :header="t('logistics.speed')" />
            <Column field="price" :header="t('common.price')" />
          </DataTable>
          <!--          </Teleport>-->
        </div>
        <Loading v-else />
      </template>
    </Card>
  </Layout>
</template>

<script setup lang="ts">
// import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable, { type DataTableRowClickEvent } from 'primevue/datatable'
import Dropdown from 'primevue/dropdown'
import FloatLabel from 'primevue/floatlabel'
import InputNumber  from 'primevue/inputnumber'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'
import Loading from '@/components/Common/Loading.vue'
import MessageBlock from '@/components/Common/MessageBlock.vue'
import { useGetData } from '@/composables/useGetData'
import { useLogistics } from '@/composables/useLogistics'
import { useMap } from '@/composables/useMap'
import { useResources } from '@/composables/useResources'
import { type BackData } from '@/types'
import type { LogisticHub, Resource, ResourceMovePayload } from '@/types/Resources/index.interface'

interface LogisticResource {
  id: number;
  name: string;
  items: Resource[];
}

const resource = ref<number>()
const from = ref<string>()
const to = ref<string>()
const amount = ref<number>()
const hub = ref<string>('')
const messageData = ref<BackData>()
const price = ref<number>(0)

// const breakpoints = useBreakpoints(breakpointsTailwind)
// const isDesktop = breakpoints.greater('md')

const { t } = useI18n()
const { getMap } = useMap()
const { getHubs } = useLogistics()
const { moveResource } = useResources()

const { data: resourcesData, onFetchResponse, isFetching } = useGetData<Resource[]>('/resource/my')
const { data: mapData } = getMap()
const { data: hubsData } = getHubs()

const groupedResources = computed<LogisticResource[]>(() => {
  if (!resourcesData.value) {
    return
  }

  return resourcesData.value.reduce((acc: LogisticResource[], item) => {
    const { resourceTypeId } = item

    const index = acc.findIndex(i => i?.id === resourceTypeId)

    if (index >= 0) {
      acc[index].items.push(item)
    } else {
      acc.push({ items: [item], name: item?.resourceType?.name, id: resourceTypeId })
    }

    return acc
  }, [] as LogisticResource[])
})

const computedResources = computed(() => {
  if (!groupedResources.value) {
    return
  }

  return groupedResources.value.map(value => {
    return {
      name: value.name,
      id: value.id
    }
  })
})

const computedFrom = computed<{cell: string; amount: number}[]>(() => {
  if (!resource.value) {
    return []
  }
  const res = groupedResources.value.find(item => item.id === resource.value)

  if (!res) {
    return []
  }

  return res.items.map(item => {
    return {
      cell: `${item.x}x${item.y}`,
      amount: item.amount
    }
  })
})

const computedTo = computed(() => {
  if (!mapData.value || !from.value) {
    return
  }

  return mapData.value.map(item => item.cellName).filter(item => item !== from.value)
})

const capacity = computed<number>(() => {
  if (!to.value || !amount.value || !resource.value) {
    return -1
  }
  const fromX = Number(from.value?.split('x')[0])
  const fromY = Number(from.value?.split('x')[1])
  const toX = Number(to.value?.split('x')[0])
  const toY = Number(to.value?.split('x')[1])

  const distance = Math.sqrt(Math.pow(toX - fromX, 2) + Math.pow(toY - fromY, 2))
  const resourceType = groupedResources.value?.find(item => item.id === resource.value)!.items[0].resourceType

  return (resourceType.weight + resourceType.volume) / 1000 * distance * amount.value
})

const hubs = computed(() => {
  if (!hubsData.value || capacity.value < 0) {
    return []
  }

  const fromX = Number(from.value?.split('x')[0])
  const fromY = Number(from.value?.split('x')[1])

  return hubsData.value.filter(hub => hub.x === fromX && hub.y === fromY).filter(hub => hub.capacity >= capacity.value).filter(hub => hub.price > 0)
})

const onRowClick = (event: DataTableRowClickEvent) => {
  hub.value = event.data.buildingId
  price.value = event.data.price
}

const rowClass = (data: LogisticHub) => {
  if (data.buildingId === hub.value) {
    return 'bg-green-100'
  }
}
const create = () => {
  if (!resource.value) {
    return
  }
  messageData.value = undefined

  const fromX = Number(from.value?.split('x')[0])
  const fromY = Number(from.value?.split('x')[1])
  const toX = Number(to.value?.split('x')[0])
  const toY = Number(to.value?.split('x')[1])

  const payload: ResourceMovePayload = {
    fromX,
    fromY,
    toX,
    toY,
    resourceTypeId: resource.value,
    amount: Number(amount.value),
    buildingId: hub.value
  }

  const { onFetchResponse: onMoveResponse, data } = moveResource(payload)
  onMoveResponse(() => {
    messageData.value = data.value
  })
}

watch(computedFrom, () => {
  from.value = computedFrom.value[0].cell
})

watch(from, () => {
  amount.value = computedFrom.value.find(item => item.cell === from.value)!.amount
})
</script>

<style scoped>

</style>