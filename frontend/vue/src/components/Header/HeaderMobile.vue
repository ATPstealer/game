<template>
  <div class="flex items-center w-full">
    <div class="flex gap-4" v-if="user?.nickName">
      <p class="font-bold text-white">
        {{ user.nickName }}
      </p>
      <p class="font-bold text-white">
        {{ moneyFormat(user?.money) }}
      </p>
    </div>
    <div v-else class="flex gap-4">
      <Button
        :label="t('account.login')"
        @click="emits('show-login')"
        class="text-white"
        severity="success"
      />
      <Button
        :label="t('account.signup')"
        @click="emits('show-signup')"
        class="text-white"
        severity="info"
      />
    </div>
    <Button
      icon="pi pi-bars"
      size="large"
      text
      @click="showSidebar = true"
      class="ml-auto text-white"
    />
    <Sidebar v-model:visible="showSidebar">
      <div class="h-full flex flex-col items-start">
        <div class="flex flex-col gap-2">
          <h3 class="underline font-bold">
            {{ t('menu.business') }}
          </h3>
          <router-link
            v-for="item in userItems"
            :key="item.label"
            :to="item.path"
            class="menu-item"
            active-class="font-bold"
          >
            {{ item.label }}
          </router-link>
          <h3 class="mt-4 underline font-bold">
            {{ t('menu.world') }}
          </h3>
          <router-link
            v-for="item in worldItems"
            :key="item.label"
            :to="item.path"
            class="menu-item"
            active-class="font-bold"
          >
            {{ item.label }}
          </router-link>
        </div>
        <div class="mt-16 flex justify-between items-center w-full">
          <Button
            :label="t('account.logout')"
            text
            severity="secondary"
            @click="emits('sign-out')"
            class="p-0"
            v-if="user?.nickName"
          />
          <LangSelect />
        </div>
      </div>
    </Sidebar>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import Sidebar from 'primevue/sidebar'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import LangSelect from '@/components/Header/LangSelect.vue'
import type { User } from '@/types'
import type { MenuItem } from '@/types/Header/index.interface'
import { moneyFormat } from '@/utils/moneyFormat'

interface Props {
  user: User | undefined;
  userItems: MenuItem[];
  worldItems: MenuItem[];
}

defineProps<Props>()

const emits = defineEmits<{
  (e: 'show-login'): void;
  (e: 'show-signup'): void;
  (e: 'sign-out'): void;
}>()

const showSidebar = ref<boolean>(false)

const { t } = useI18n()

</script>

<style scoped>
.menu-item {
  @apply no-underline text-gray-800 visited:text-gray-800;
}
</style>
