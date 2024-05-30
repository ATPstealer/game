<template>
  <PipelinesTemplate :headers="headers">
    <template #items>
      <div v-if="blueprints?.length" class="flex gap-2 flex-wrap">
        <span
          v-for="blueprint in blueprints"
          :key="blueprint.id"
          class="pipeline-item pipeline-item-hover"
          :class="{'bg-amber-100' : blueprint.id === chosen?.id }"
          @click="chosen = blueprint"
        >
          {{ t(`blueprints.types.${blueprint?.name.toLowerCase()}`) }}
        </span>
      </div>
    </template>
    <template #first-column>
      <div v-if="computedData?.usedBlueprints?.length" class="flex flex-col gap-2">
        <div v-for="item in computedData.usedBlueprints" :key="item.id">
          <span class="pipeline-item pipeline-item-hover" @click="chosen = item">{{ t(`buildings.types.${item.name.toLowerCase()}`) }}</span>
        </div>
      </div>
    </template>
    <template #second-column>
      <div v-if="chosen" class="flex flex-col gap-2">
        <span class="text-xl underline text-center">{{ t(`buildings.types.${chosen?.name.toLowerCase()}`) }}</span>
      </div>
    </template>
    <template #third-column>
      <div
        v-if="computedData?.prodBlueprints?.length"
        class="flex flex-col gap-2"
        :class="{'pipeline-column-grid': computedData?.prodBlueprints?.length > 10}"
      >
        <div v-for="item in computedData.prodBlueprints" :key="item.id">
          <span class="pipeline-item pipeline-item-hover" @click="chosen = item">{{ t(`buildings.types.${item.name.toLowerCase()}`) }}</span>
        </div>
      </div>
    </template>
  </PipelinesTemplate>
</template>

<script setup lang="ts">
import { computed, ref, toRef } from 'vue'
import { useI18n } from 'vue-i18n'
import PipelinesTemplate from '@/components/Pipelines/PipelinesTemplate.vue'
import type { Blueprint } from '@/types/Buildings/index.interface'

interface Props {
  blueprints: Blueprint[];
}

interface Pipeline {
  usedBlueprints: Blueprint[];
  prodBlueprints: Blueprint[];
}

const props = defineProps<Props>()
const blueprints = toRef(props, 'blueprints')

const { t } = useI18n()

const chosen = ref<Blueprint>(blueprints?.value?.[0] || {} as Blueprint)
// TODO: добавить переводы
const headers = ['Нужны ресурсы из', t('common.blueprint'), 'Производит ресурсы лдя']

const computedData = computed<Pipeline>(() => {
  if (!chosen.value) {
    return {} as Pipeline
  }

  const usedResources = chosen.value?.usedResources.map(resource => resource.resourceId)
  const prodResources = chosen.value?.producedResources.map(resource => resource.resourceId)

  return {
    usedBlueprints: blueprints.value.filter(bp => bp.producedResources.some(resource => usedResources?.includes(resource.resourceId))),
    prodBlueprints: blueprints.value.filter(bp => bp.usedResources.some(resource => prodResources?.includes(resource.resourceId)))
  }
})
</script>

<style scoped>

</style>