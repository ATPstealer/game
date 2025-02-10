<template>
  <Layout :show-options-prop="false">
    <div v-if="!isFetching && equipments?.length">
      <EquipmentList
        :equipments="equipments"
      />
    </div>
    <template #help>
      Equipment help
    </template>
  </Layout>
</template>

<script setup lang="ts">
import { computed, unref } from 'vue'
import Layout from '@/components/Common/Layout.vue'
import EquipmentList from '@/components/Equipment/EquipmentList.vue'
import { useGetEquipmentMy } from '@/gen'

const { data, suspense, isFetching } = useGetEquipmentMy()
await suspense()
const equipments = computed(() => unref(data)?.data)

</script>
