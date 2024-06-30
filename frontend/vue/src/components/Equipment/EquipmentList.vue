<template>
  <DataTable
    size="small"
    striped-rows
    :value="equipments"
  >
    <Column
      :header="t(`equipment.columns.name`)"
    >
      <template #body="{data}: {data: Equipment}">
        {{ t(`resources.types.${data.resourceType.name.toLowerCase()}`) }}
      </template>
    </Column>
    <Column
      :header="t(`equipment.columns.effect`)"
    >
      <template #body="{data}: {data: Equipment}">
        {{ t(`equipment.effect.${data.equipmentType.effectId.toString()}`) }}
      </template>
    </Column>
    <Column
      :header="t(`common.blueprint`)"
    >
      <template #body="{data}: {data: Equipment}">
        <span v-if="!data.equipmentType.blueprintIds.length">Все</span>
        <span
          v-else
          class="text-indigo-500 cursor-pointer"
          @click="toggleOP($event, data.equipmentType.blueprintIds)"
        >
          Посмотреть
        </span>
        <OverlayPanel ref="popover">
          <div
            v-for="bp in bps"
            :key="bp.id"
            class="flex flex-col gap-8"
          >
            <router-link
              v-slot="{href, navigate}"
              custom
              :to="{name: 'Pipelines', query: {selected: 'blueprint',id: bp.id, }}"
            >
              <a
                :href="href"
                target="_blank"
                @click="navigate"
              >
                {{ bp.name }}
              </a>
            </router-link>
          </div>
        </OverlayPanel>
      </template>
    </Column>
    <Column
      :header="t(`common.value`)"
    >
      <template #body="{data}: {data: Equipment}">
        {{ data.equipmentType.value }}
      </template>
    </Column>
    <Column
      :header="t(`map.square`)"
    >
      <template #body="{data}: {data: Equipment}">
        {{ data.equipmentType.square }}
      </template>
    </Column>
    <Column
      :header="t(`equipment.columns.durability`)"
    >
      <template #body="{data}: {data: Equipment}">
        {{ data.equipmentType.durability }}
      </template>
    </Column>
    <Column
      field="amount"
      :header="t(`common.amount`)"
    />
    <Column
      :header="t(`map.cell`)"
    >
      <template #body="{data}: {data: Resource}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>
  </DataTable>
</template>

<script setup lang="ts">
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import OverlayPanel from 'primevue/overlaypanel'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useGetData } from '@/composables/useGetData'
import type { Blueprint } from '@/types/Buildings/index.interface'
import type { Equipment } from '@/types/Equipment/index.interface'
import { Resource } from '@/types/Resources/index.interface'

interface Props {
  equipments: Equipment[];
}
const props = defineProps<Props>()

const { t } = useI18n()

const popover = ref()
const bps = ref<Blueprint[]>([])

const toggleOP = (event: any, blueprintsIds: number[]) => {
  popover.value.toggle(event)
  bps.value = blueprints.value.filter(item => blueprintsIds.includes(item.id))
}

const { data: blueprints } = useGetData('/building/blueprints')

</script>