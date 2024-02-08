<template>
  <Layout :show-options-prop="false">
    <Orders />
    <template #help>
      {{ t(`orders.help`) }}
    </template>
  </Layout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'
import Orders from '@/components/Orders/Orders.vue'

const { t } = useI18n()

</script>
