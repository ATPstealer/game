<template>
  <ResourcesList
    :resources="resources"
    :execute="execute"
    v-if="resources?.length"
  />
  <div v-else>
    тут рыбы нет
  </div>
</template>

<script setup lang="ts">
import isEmpty from 'lodash/isEmpty'
import { ref, watch } from 'vue'
import ResourcesList from '@/components/Resources/ResourcesList.vue'
import { useGetData } from '@/composables/useGetData'
import type { Building } from '@/types/Buildings/index.interface'
import type { Resource } from '@/types/Resources/index.interface'

interface Props {
  building: Building;
}

const props = defineProps<Props>()

const resources = ref<Resource[]>([])

const { data, onFetchResponse, execute } = useGetData<Resource[]>(`/resource/my?x=${props.building.x}&y=${props.building.y}`, false)

onFetchResponse(() => {
  resources.value = data.value
})

watch(() => props.building, () => {
  if (!isEmpty(props.building)) {
    execute()
  }
}, {
  immediate: true
})
</script>

<style scoped>

</style>