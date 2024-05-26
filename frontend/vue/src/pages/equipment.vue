<template>
  <Layout :show-options-prop="false">
    Equipment
    <template #help>
      Equipment
    </template>
  </Layout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'

const { t } = useI18n()

</script>
