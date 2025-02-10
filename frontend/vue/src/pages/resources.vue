<template>
  <Layout :show-options-prop="false">
    <div v-if="!isFetching && data?.length">
      <Resources
        :execute="refetch"
        :resources="data.filter(resource => resource?.resourceType?.name)"
      />
    </div>
    <Loading v-else />
    <template #help>
      {{ t(`resources.help`) }}
    </template>
  </Layout>
</template>

<script setup lang="ts">
import { computed, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'
import Loading from '@/components/Common/Loading.vue'
import Resources from '@/components/Resources/ResourcesList.vue'
import { useGetResourceMy } from '@/gen'

const { t } = useI18n()

const { data: resourcesQuery, isFetching, suspense, refetch } = useGetResourceMy()
await suspense()
const data = computed(() => unref(resourcesQuery)?.data)

</script>
