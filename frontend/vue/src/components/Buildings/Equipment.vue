<template>
  <div class="flex flex-col gap-4">
    <Divider />
    <p class="font-bold capitalize">
      {{ t(`common.equipment`) }}
    </p>
    <Button
      class="w-max"
      :label="t('equipment.title')"
      size="small"
      @click="router.push({name: 'EquipmentId', params: {id: building?._id}})"
    />
    <div v-if="building?.equipmentEffect?.length && blueprints?.length">
      <p class="font-bold mb-4">
        Эффекты оборудования:
      </p>
      <div class="flex flex-col gap-2">
        <Tag
          v-for="effect in building.equipmentEffect"
          :key="effect.effectId"
          v-tooltip.top="`Влияние: ${effect.value.toFixed(3)} на ${effect.blueprintId === 0 ? 'все чертежи' : blueprints.find(bp => bp.id === effect.blueprintId)?.name}`"
          severity="info"
          :value="t(`equipment.effect.${effect.effectId.toString()}`)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Divider from 'primevue/divider'
import Tag from 'primevue/tag'
import { computed, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import { type BuildingWithData, useGetBuildingBlueprints } from '@/gen'
import router from '@/router'

interface Props {
  building: BuildingWithData | undefined;
}

defineProps<Props>()

const { t } = useI18n()

const { data: blueprintsQuery, suspense: awaitBlueprints } = useGetBuildingBlueprints()
await awaitBlueprints()
const blueprints = computed(() => unref(blueprintsQuery)?.data)
</script>