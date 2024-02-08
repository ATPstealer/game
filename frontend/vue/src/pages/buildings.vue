<template>
  <Layout>
    <template #options>
      <router-link
        class="link"
        :to="{name: 'ConstructBuilding', query: {x: 0, y: 0}}"
      >
        <Button :label="t(`buildings.construct new building`)" />
      </router-link>
    </template>
    <Buildings />
  </Layout>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import { useI18n } from 'vue-i18n'
import Buildings from '@/components/Buildings/Buildings.vue'
import Layout from '@/components/Common/Layout.vue'

const { t } = useI18n()
</script>

<style scoped>

</style>
