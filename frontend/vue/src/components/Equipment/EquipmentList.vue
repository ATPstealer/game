<template>
  <DataTable
      :value="equipments"
      size="small"
      striped-rows
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
        :header="t(`map.cell`)"
    >
      <template #body="{data}: {data: Resource}">
        {{ data.x }}x{{ data.y }}
      </template>
    </Column>

  </DataTable>
</template>


<script setup lang="ts">
import type {Equipment} from "@/types/Equipment/index.interface";
import Column from "primevue/column";
import DataTable from "primevue/datatable";
import {useI18n} from "vue-i18n";
import {Resource} from "@/types/Resources/index.interface";


interface Props {
  equipments: Equipment[];
}
const { t } = useI18n()
const props = defineProps<Props>()
const equipments = props.equipments


</script>