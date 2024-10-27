<template>
  <Layout>
    <template #options>
      <div class="flex flex-col gap-8">
        <Dropdown
          v-model="buildingType"
          option-label="title"
          :options="buildingTypes"
          @change="(event: DropdownChangeEvent) => setParams({key: 'buildingTypeId', value: event.value.id})"
        >
          <template #option="{option}: {option: BuildingType}">
            {{ getTranslation({parent: 'buildings.types', child: option.title}) }}
          </template>
          <template #value="{value}: {value: BuildingType}">
            {{ getTranslation({parent: 'buildings.types', child: value.title}) }}
          </template>
        </Dropdown>
        <div class="flex gap-4">
          <FloatLabel class="flex-1">
            <Dropdown
              v-model="x"
              class="w-full"
              input-id="x"
              :options="xOptions"
              placeholder="X"
              @change="(event: DropdownChangeEvent) => setParams({key: 'x', value: event.value})"
            >
              <template #value="{value}: {value: number}">
                {{ value === undefined ? 'X' : value }}
              </template>
            </Dropdown>
            <label for="x">X</label>
          </FloatLabel>
          <FloatLabel class="flex-1">
            <Dropdown
              v-model="y"
              class="w-full"
              input-id="y"
              :options="yOptions"
              placeholder="Y"
              @change="(event: DropdownChangeEvent) => setParams({key: 'y', value: event.value})"
            >
              <template #value="{value}: {value: number}">
                {{ value === undefined ? 'Y' : value }}
              </template>
            </Dropdown>
            <label for="y">Y</label>
          </FloatLabel>
        </div>
        <AutoComplete
          v-model="userSearch"
          :delay="1000"
          dropdown
          :placeholder="t('search-building.chooseUser')"
          :suggestions="users"
          @complete="searchUsers"
          @item-select="event => setParams({key: 'nickName', value: event.value})"
        />
        <Button :label="t('common.reset')" @click="clearParams" />
      </div>
    </template>
    <BuildingSearch :search-params="params" />
  </Layout>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import AutoComplete from 'primevue/autocomplete'
import Button from 'primevue/button'
import Dropdown, { type DropdownChangeEvent } from 'primevue/dropdown'
import FloatLabel from 'primevue/floatlabel'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { BuildingType, FindBuildingParams } from '@/api'
import { getBuildingTypesOptions, getSettingsOptions } from '@/api/@tanstack/vue-query.gen'
import BuildingSearch from '@/components/Buildings/SearchBuilding/BuildingSearch.vue'
import Layout from '@/components/Common/Layout.vue'
import { useGetData } from '@/composables/useGetData'
import type { Coords } from '@/types'
import type { SearchBuildingParams } from '@/types/Buildings/index.interface'
import { getTranslation } from '@/utils/getTranslation'

const buildingType = ref<BuildingType>({ id: 0, title: 'All' } as BuildingType)
const buildingTypes = ref<BuildingType[]>([{ id: 0, title: 'All' } as BuildingType])
const x = ref<number>()
const y = ref<number>()
const userSearch = ref<string>('')
const users = ref<string[]>([])
const coords = ref<Record<Coords, number>>()

const { t } = useI18n()

const { suspense: gotBuildingTypes } = useQuery({
  ...getBuildingTypesOptions(),
  select: (data: any) => {
    buildingTypes.value = [...buildingTypes.value, ...data.data]
  }
})

const { suspense: gotCoords } = useQuery({
  ...getSettingsOptions(),
  select: (data: any) => {
    coords.value = data.data
  }
})

const searchUsers = () => {
  const { data, onFetchResponse } = useGetData<string[]>(`/data/users_by_prefix?prefix=${userSearch.value}`)
  onFetchResponse(() => {
    users.value = data.value
  })
}

const params = ref<FindBuildingParams>({
  limit: 500000000
})

const setParams = <Key extends keyof SearchBuildingParams>({ key, value }: {key: Key; value: SearchBuildingParams[Key]}) => {
  params.value[key] = value

  if (!buildingType.value.id) {
    delete params.value.buildingTypeId
  }
}

const clearParams = () => {
  params.value = {
    limit: 500000000
  }
  userSearch.value = ''
  x.value = undefined
  y.value = undefined
}

const xOptions = computed(() => {
  const minX = coords.value!.mapMinX
  const maxX = coords.value!.mapMaxX

  return getValues(minX, maxX)
})

const yOptions = computed(() => {
  const minY = coords.value!.mapMinY
  const maxY = coords.value!.mapMaxY

  return getValues(minY, maxY)
})

const getValues = (min: number, max: number) => {
  const values = []
  for (let i = min; i <= max; i++) {
    values.push(i)
  }

  return values
}

await gotBuildingTypes()
await gotCoords()
</script>

<style scoped>
</style>
