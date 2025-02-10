<template>
  <Layout>
    <template #options>
      <div v-for="option in options" :key="option">
        <RadioButton
          v-model="selected"
          :name="option"
          :value="option"
        />
        <label class="ml-2" :for="option">{{ t(`common.${option}`) }}</label>
      </div>
    </template>
    <template v-if="resourceTypes?.length && blueprints?.length && buildingTypes?.length">
      <ResourcesPipeline
        v-if="selected === 'resource'"
        :blueprints="blueprints"
        :building-types="buildingTypes"
        :resource-types="resourceTypes"
      />
      <BuildingsPipeline
        v-if="selected === 'building'"
        :blueprints="blueprints"
        :building-types="buildingTypes"
        :resource-types="resourceTypes"
      />
      <BlueprintsPipeline
        v-if="selected === 'blueprint'"
        :blueprints="blueprints"
        :selected-id="+route?.query.id"
      />
    </template>
  </Layout>
</template>

<script setup lang="ts">
import RadioButton from 'primevue/radiobutton'
import { computed, ref, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import Layout from '@/components/Common/Layout.vue'
import BlueprintsPipeline from '@/components/Pipelines/BlueprintsPipeline.vue'
import BuildingsPipeline from '@/components/Pipelines/BuildingsPipeline.vue'
import { options } from '@/components/Pipelines/constants'
import ResourcesPipeline from '@/components/Pipelines/ResourcesPipeline.vue'
import { useGetBuildingBlueprints, useGetBuildingTypes, useGetResourceTypes } from '@/gen'

const route = useRoute()
const { t } = useI18n()

const selected = ref(route.query.selected || 'resource')

const { data: resourceTypesQuery, suspense: awaitResourceTypes } = useGetResourceTypes()
await awaitResourceTypes()
const resourceTypes = computed(() => unref(resourceTypesQuery)?.data)

const { data: blueprintsTypesQuery, suspense: awaitBlueprints } = useGetBuildingBlueprints()
await awaitBlueprints()
const blueprints = computed(() => unref(blueprintsTypesQuery)?.data)

const { data: buildingTypesQuery, suspense: awaitBuildingTypes } = useGetBuildingTypes()
await awaitBuildingTypes()
const buildingTypes = computed(() => unref(buildingTypesQuery)?.data)

</script>

<style>

</style>