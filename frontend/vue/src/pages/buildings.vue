<template>
  <Layout>
    <template #options>
      <router-link
          class="text-blue-500 hover:text-blue-700"
          :to="{name: 'ConstructBuilding', query: {x: 0, y: 0}}"
      >
        {{ t(`buildings.construct new building`) }}
      </router-link>
    </template>
    <Buildings />
  </Layout>
</template>

<script setup lang="ts">
import Layout from "@/components/Common/Layout.vue";
import Buildings from "@/components/Buildings/Buildings.vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()
</script>

<style scoped>

</style>
