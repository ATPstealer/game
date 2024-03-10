<template>
  <Layout>
    <template #options>
      <Dropdown
        :options="buildingTypes"
        option-label="title"
        v-model="buildingType"
        @change="event => setParams({key: 'buildingTypeId', value: event.value.id})"
      >
        <template #option="{option}: {option: BuildingType}">
          {{ getTranslation({parent: 'buildings.types', child: option.title}) }}
        </template>
        <template #value="{value}: {value: BuildingType}">
          {{ getTranslation({parent: 'buildings.types', child: value.title}) }}
        </template>
      </Dropdown>
      <div class="slider">
        <span>X: {{ x }}</span>
        <div class="flex items-center gap-4">
          <span class="text-sm">{{ coords?.mapMinX }}</span>
          <Slider
            v-model="x"
            :max="coords?.mapMaxX"
            :min="coords?.mapMinX"
            @change="value => x = value"
            class="w-[150px]"
            @slideend="event => setParams({key: 'x', value: event.value})"
          />
          <span class="text-sm">{{ coords?.mapMaxX }}</span>
        </div>
      </div>
      <div class="slider">
        <span>Y: {{ y }}</span>
        <div class="flex items-center gap-4">
          <span class="text-sm">{{ coords?.mapMinY }}</span>
          <Slider
            v-model="y"
            :max="coords?.mapMaxY"
            :min="coords?.mapMinY"
            @change="value => y = value"
            class="w-[150px]"
            @slideend="event => setParams({key: 'y', value: event.value})"
          />
          <span class="text-sm">{{ coords?.mapMaxY }}</span>
        </div>
      </div>
      <AutoComplete
        v-model="userSearch"
        :suggestions="users"
        :delay="1000"
        @complete="searchUsers"
        @item-select="event => setParams({key: 'nickName', value: event.value})"
        dropdown
        :placeholder="t('search-building.chooseUser')"
      />

      <Button :label="t('common.reset')" @click="clearParams" />
    </template>
    <BuildingSearch :search-params="params" />
  </Layout>
</template>

<script setup lang="ts">
import AutoComplete from 'primevue/autocomplete'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import Slider from 'primevue/slider'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import BuildingSearch from '@/components/Buildings/SearchBuilding/BuildingSearch.vue'
import Layout from '@/components/Common/Layout.vue'
import { useGetData } from '@/composables/useGetData'
import type { Coords } from '@/types'
import type { BuildingType, SearchBuildingParams } from '@/types/Buildings/index.interface'
import { getTranslation } from '@/utils/getTranslation'

const buildingType = ref<BuildingType>({ id: 0, title: 'All' })
const buildingTypes = ref<BuildingType[]>([{ id: 0, title: 'All' }])
const x = ref<number>()
const y = ref<number>()
const userSearch = ref<string>('')
const users = ref<string[]>([])

const { t } = useI18n()

const { data, onFetchResponse } = useGetData<BuildingType[]>('/building/types')
onFetchResponse(() => {
  buildingTypes.value = [...buildingTypes.value, ...data.value]
})

const { data: coords } = useGetData<Record<Coords, number>>('/settings')

const searchUsers = () => {
  const { data, onFetchResponse } = useGetData<string[]>(`/data/users_by_prefix?prefix=${userSearch.value}`)
  onFetchResponse(() => {
    users.value = data.value
  })
}

const params = ref<SearchBuildingParams>({
  limit: 500000000 // TODO: Илья сделай плз пагинатор
})

const setParams = ({ key, value }: {key: string; value: number}) => {
  params.value[key] = value

  if (!buildingType.value.id) {
    delete params.value.buildingTypeId
  }
}

const clearParams = () => {
  params.value = {
    limit: -1
  }
  userSearch.value = ''
}

</script>

<style scoped>
.slider {
  @apply w-full flex items-center justify-between;
}
</style>
