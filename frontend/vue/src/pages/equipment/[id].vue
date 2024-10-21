<template>
  <Card class="max-w-screen-xl m-auto">
    <template #content>
      <div class="flex gap-4 items-center">
        <Button
          class="w-5 h-9 shadow-md"
          icon="pi pi-arrow-left"
          rounded
          size="small"
          text
          @click="returnToBuilding"
        />
        <h2>{{ t('equipment.title') }}</h2>
      </div>
      <div class="flex flex-col xl:flex-row gap-12 xl:gap-24 items-center mt-8 mb-12">
        <div class="flex flex-col gap-4">
          <p>{{ t('common.total') }} {{ t('common.square').toLocaleLowerCase() }}: {{ building.square }} / {{ t('equipment.available') }}: {{ availableSpace }}</p>
          <Chart
            class="w-full h-full xl:w-[500px] xl:h-[500px]"
            :data="chartData"
            :options="chartOptions"
            type="doughnut"
          />
        </div>
        <InstalledEquipment
          :building="building"
          :equipment-types="equipmentTypes"
          @update-equipment="() => refetchBuilding()"
        />
      </div>
      <EquipmentToInstall :building="building" @update-equipment="() => refetchBuilding()" />
    </template>
  </Card>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Chart from 'primevue/chart'
import type RouteNamedMap from 'typed-router'
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import type { BuildingWithData, EquipmentType } from '@/api'
import {
  getBuildingMyOptions,
  getEquipmentTypesOptions
} from '@/api/@tanstack/vue-query.gen'
import { colors } from '@/components/Equipment/constants'
import EquipmentToInstall from '@/components/Equipment/EquipmentToInstall.vue'
import InstalledEquipment from '@/components/Equipment/InstalledEquipment.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const building = ref<BuildingWithData>({} as BuildingWithData)
const equipmentTypes = ref<EquipmentType[]>([])
const chartOptions = ref()

const { suspense: gotBuilding, refetch: refetchBuilding } = useQuery({
  ...getBuildingMyOptions({ query: { _id: route.params.id } }),
  select: (data: any) => {
    building.value = data.data[0]

    return data
  }
})

const { suspense: gotEquipmentTypes } = useQuery({
  ...getEquipmentTypesOptions(),
  select: (data: any) => {
    equipmentTypes.value = data.data
  }
})

const chartData = computed(() => {
  const allEquipment = building.value.equipment
  let initialData = [] as { label: string; weight: number; amount: number }[]
  if (!allEquipment?.length) {
    initialData = []
  }

  else {
    initialData = allEquipment.map(item => {
      const eq = equipmentTypes.value.find(type => type.id === item.equipmentTypeId)!

      return {
        label: eq.name,
        weight: item.amount * eq.square,
        amount: item.amount
      }
    })
  }

  const documentStyle = getComputedStyle(document.body)

  const fullWeight = initialData.map(item => item.weight).reduce((a, b) => a + b, 0)
  const labels = initialData.map(item => `${item.label} ${item.amount}${t('common.pieces')}/${item.weight}`)
  const data = initialData.map(item => item.weight)

  const backgroundColor = []
  data.forEach((item, index) => backgroundColor.push(documentStyle.getPropertyValue(index > colors.length ? `--${colors[index-11]}-500` : `--${colors[index]}-500`)))

  const hoverBackgroundColor = []
  data.forEach((item, index) => hoverBackgroundColor.push(documentStyle.getPropertyValue(`--${colors[index]}-400`)))

  if (fullWeight !== building.value?.square) {
    labels.push(t('equipment.available'))
    data.push(building.value?.square - fullWeight)
    backgroundColor.push(documentStyle.getPropertyValue('--gray-300'))
    hoverBackgroundColor.push(documentStyle.getPropertyValue('--gray-200'))
  }

  return {
    labels,
    datasets: [
      {
        data,
        backgroundColor,
        hoverBackgroundColor
      }
    ]
  }
})

const availableSpace = computed(() => {
  if (!building.value.equipment) {
    return
  }

  return building.value.square - building.value.equipment.map(item => item.amount * equipmentTypes.value.find(eq => eq.id === item.equipmentTypeId)!.square).reduce((a, b) => a + b, 0)
})

const setChartOptions = () => {
  const documentStyle = getComputedStyle(document.documentElement)
  const textColor = documentStyle.getPropertyValue('--text-color')

  return {
    plugins: {
      legend: {
        labels: {
          cutout: '60%',
          color: textColor
        },
        position: 'bottom'
      }
    }
  }
}

const returnToBuilding = () => {
  const routeName = `Buildings${building.value.buildingType.buildingGroup}NameId` as unknown as keyof typeof RouteNamedMap | undefined
  router.push({ name: routeName, params: { id: building.value._id, name: building.value.buildingType.title.toLowerCase() } })
}

onMounted(() => {
  chartOptions.value = setChartOptions()
})

await gotBuilding()
await gotEquipmentTypes()
</script>

<style scoped>

</style>