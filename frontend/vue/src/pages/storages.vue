<template>
  <Layout :show-options-prop="false">
    <Storages />
    <template #help>
      {{ t(`storages.help`) }}
    </template>
  </Layout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'
import Storages from '@/components/Storages/Storages.vue'

const { t } = useI18n()

</script>
<style scoped>

</style>
