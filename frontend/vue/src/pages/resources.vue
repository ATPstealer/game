<template>
  <Layout :show-options-prop="false">
    <div v-if="!isFetching && data?.length">
      <Resources
        :execute="execute"
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
import { useI18n } from 'vue-i18n'
import Layout from '@/components/Common/Layout.vue'
import Loading from '@/components/Common/Loading.vue'
import Resources from '@/components/Resources/ResourcesList.vue'
import { useGetData } from '@/composables/useGetData'
import type { Resource } from '@/types/Resources/index.interface'

const { t } = useI18n()

const { data, isFetching, execute } = useGetData<Resource[]>('/resource/my')

</script>
