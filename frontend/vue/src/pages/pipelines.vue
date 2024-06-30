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
    <template v-if="resourceTypes?.length && blueprints?.length && buildingsTypes?.length">
      <ResourcesPipeline
        v-if="selected === 'resource'"
        :blueprints="blueprints"
        :buildings-types="buildingsTypes"
        :resource-types="resourceTypes"
      />
      <BuildingsPipeline
        v-if="selected === 'building'"
        :blueprints="blueprints"
        :buildings-types="buildingsTypes"
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
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import Layout from '@/components/Common/Layout.vue'
import BlueprintsPipeline from '@/components/Pipelines/BlueprintsPipeline.vue'
import BuildingsPipeline from '@/components/Pipelines/BuildingsPipeline.vue'
import { options } from '@/components/Pipelines/constants'
import ResourcesPipeline from '@/components/Pipelines/ResourcesPipeline.vue'
import { useGetData } from '@/composables/useGetData'
import { Blueprint, Building } from '@/types/Buildings/index.interface'
import { Resource } from '@/types/Resources/index.interface'

const route = useRoute()

const { data: resourceTypes } = useGetData<Resource[]>('/resource/types')
const { data: blueprints } = useGetData<Blueprint[]>('/building/blueprints')
const { data: buildingsTypes } = useGetData<Building[]>('/building/types')

const selected = ref(route.query.selected || 'resource')

const { t } = useI18n()
</script>

<style>

</style>