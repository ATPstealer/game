<template>
  <Header />
  <div class="px-0 md:px-16 py-4 relative h-full">
    <ConfirmPopup
      :pt="{
        message: {
          class: 'max-w-[250px]'
        }
      }"
    />
    <Breadcrumb
      v-if="breadcrumbs.length"
      :home="home"
      :model="breadcrumbs"
      :pt="{
        root: {
          class: 'md:w-[500px] md:mb-4 md:-ml-12'
        }
      }"
    >
      <template #item="{ item, props }">
        <router-link
          v-if="item.route"
          v-slot="{ href, navigate }"
          custom
          :to="item.route"
        >
          <a
            :href="href"
            v-bind="props.action"
            @click="navigate"
          >
            <span :class="[item.icon, 'text-color']" />
            <span class="text-primary font-semibold">{{ item.label }}</span>
          </a>
        </router-link>
        <span v-else class="text-color">{{ item.label }}</span>
      </template>
    </Breadcrumb>
    <suspense>
      <RouterView />
    </suspense>
    <Toast />
  </div>
  <div class="footer">
    футер
  </div>
</template>

<script setup lang="ts">
import Breadcrumb from 'primevue/breadcrumb'
import ConfirmPopup from 'primevue/confirmpopup'
import Toast from 'primevue/toast'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterView, useRoute } from 'vue-router'
import Header from '@/components/Header/Header.vue'

const route = useRoute()
const { t } = useI18n()

const home = ref({
  icon: 'pi pi-home',
  route: '/'
})

const breadcrumbs = computed(() => {
  const { path } = route
  const items = path.split('/')
  items.shift()
  const pageName = items.pop()?.split('-')[0]
  items.pop()

  const page = items[0]

  const result: {label: string; route?: string}[] = items.map(item => {
    return {
      label: t(`${item}.title`),
      route: `/${item}`
    }
  })

  if (items.length) {
    result.push({
      label: t(`${page}.types.${decodeURI(pageName)}`)
    })
  }

  return result
})
</script>

<style scoped>
.footer {
  @apply mt-auto;
}
</style>
