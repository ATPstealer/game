<template>
  <Layout :show-options-prop="false">
    <Button :label="t('logistics.send')" @click="$router.push({ name: 'LogisticsCreate' })" />
    <Logistics />
    <template #help>
      {{ t(`logistics.help`) }}
    </template>
  </Layout>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'
import Logistics from '@/components/Logistics/Logistics.vue'

const { t } = useI18n()
</script>
