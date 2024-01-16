<template>
  <Layout :show-options-prop="false">
    <Resources />
    <template #help>
      {{ t(`resources.help`) }}
    </template>
  </Layout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'
import Resources from '@/components/Resources/Resources.vue'

const { t } = useI18n()
</script>
