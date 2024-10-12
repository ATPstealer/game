<template>
  <Layout :show-options-prop="false">
    <Characteristics />
    {{ messageData }}
    <pre v-if="!isPending">{{ data }} </pre>
    <p v-else>
      ...loading
    </p>
  </Layout>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import {getBuildingTypes} from '@/api'
import Characteristics from '@/components/Characteristics/Characteristics.vue'
import Layout from '@/components/Common/Layout.vue'

const messageData = ref()

const { data, isPending  } = useQuery({
  queryKey: ['buildingTypes'],
  queryFn: getBuildingTypes
})

</script>

<style scoped>

</style>
