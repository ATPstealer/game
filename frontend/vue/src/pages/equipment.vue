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
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'
import EquipmentList from '@/components/Equipment/EquipmentList.vue'
import { useGetData } from '@/composables/useGetData'
import type {Equipment} from "@/types/Equipment/index.interface";

const { t } = useI18n()
const { data: equipments, isFetching, execute } = useGetData<Equipment[]>('/equipment/my')

</script>
