<template>
  <Layout>
    <Resources />
    <template #help>
      {{ t(`resources.help`) }}
    </template>
  </Layout>
</template>

<script setup lang="ts">
import Layout from "@/components/Common/Layout.vue";
import Resources from '@/components/Resources/Resources.vue'
import {useI18n} from "vue-i18n";
const { t } = useI18n()
</script>
