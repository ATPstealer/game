<template>
  <Layout class="mt-8" :show-options-prop="false">
    <Card class="max-w-screen-xl mx-auto mt-16">
      <template #content>
        <div v-if="!isFetching" class="flex flex-col gap-4">
          <MessageBlock
            v-if="messageData"
            class="mb-8"
            :code="messageData.code"
            @close-message="messageData = undefined"
          />
          <div class="flex gap-4 items-center" :class="{'mt-28': !messageData}">
            <FloatLabel>
              <Dropdown
                id="resource"
                v-model="resource"
                class="max-w-[200px] min-w-[200px]"
                :option-label="e => t(`resources.types.${e.name?.toLowerCase()}`)"
                option-value="id"
                :options="computedResources"
                placeholder="Выберите ресурс"
              />
              <label for="resource">Выберите ресурс</label>
            </FloatLabel>
            <FloatLabel>
              <Dropdown
                id="from"
                v-model="from"
                class="max-w-[200px] min-w-[200px]"
                :disabled="!resource || computedFrom.length === 1"
                :option-label="e => `${e.cell} (${e.amount})`"
                :option-value="e => e.cell"
                :options="computedFrom"
                placeholder="Откуда"
              />
              <label for="from">Откуда</label>
            </FloatLabel>
            <FloatLabel>
              <Dropdown
                id="to"
                v-model="to"
                class="max-w-[150px] min-w-[150px]"
                :disabled="!from"
                filter
                :options="computedTo"
                placeholder="Куда"
              />
              <label for="to">Куда</label>
            </FloatLabel>
            <FloatLabel>
              <InputNumber
                id="amount"
                v-model="amount"
                :disabled="!to"
                :max="computedFrom?.find(item => item.cell === from)?.amount"
                show-buttons
                @input="e => amount = +e.value"
              />
              <label for="amount">Количество</label>
            </FloatLabel>
            <div v-tooltip="!hub ? 'Выберите хаб' : ''">
              <Button
                :disabled="!hub"
                label="Отправить груз"
                severity="info"
                type="submit"
                @click="create"
              />
            </div>
          </div>
          <div>
            <span>Необходимая вместимость: </span> <span v-if="capacity >= 0">{{ capacity }}</span>
          </div>

          <DataTable
            :row-class="rowClass"
            selection-mode="single"
            :value="hubs"
            @row-click="onRowClick"
          >
            <template #empty>
              <span v-if="capacity >= 0">В данной ячейке нет хабов с нужной вместимостью, попробуйте выбрать другую ячейку или уменьшить количество груза</span>
              <span v-else>Сначала заполните другие шаги</span>
            </template>
            <Column field="capacity" header="Вместимость" />
            <Column field="speed" header="Скорость" />
            <Column field="price" header="Стоимость" />
          </DataTable>
        </div>
        <Loading v-else />
      </template>
    </Card>
  </Layout>
</template>

<script setup lang="ts">
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

  return hubsData.value.filter(hub => hub.x === fromX && hub.y === fromY).filter(hub => hub.capacity >= capacity.value)
})

const onRowClick = (event: DataTableRowClickEvent) => {
  hub.value = event.data.buildingId
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