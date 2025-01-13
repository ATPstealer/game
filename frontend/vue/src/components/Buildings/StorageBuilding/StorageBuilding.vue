<template>
  <ResourcesList
    v-if="resources?.length"
    :execute="refetch"
    :resources="resources"
  />
  <div v-else>
    тут рыбы нет
  </div>
</template>

<script setup lang="ts">
import { computed, unref } from 'vue'
import ResourcesList from '@/components/Resources/ResourcesList.vue'
import { type BuildingWithData, useGetResourceMy } from '@/gen'

interface Props {
  building: BuildingWithData | undefined;
}

const props = defineProps<Props>()

const { data: resourcesQuery, refetch, suspense } = useGetResourceMy( { x: props.building?.x, y: props.building?.y } )
await suspense()
const resources = computed(() => unref(resourcesQuery)?.data || [])

// watch(() => props.building, () => {
//   if (!isEmpty(props.building)) {
//     refetch()
//   }
// }, {
//   immediate: true
// })
</script>

<style scoped>

</style>