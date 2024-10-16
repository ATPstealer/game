<template>
  <div class="flex items-center w-full">
    <div v-if="user?.nickName" class="flex gap-4">
      <p class="font-bold text-white">
        {{ user.nickName }}
      </p>
      <p class="font-bold text-white">
        {{ moneyFormat(user?.money) }}
      </p>
    </div>
    <!--    <div v-else class="flex gap-4">-->
    <!--      <Button-->
    <!--        class="text-white"-->
    <!--        :label="t('account.login')"-->
    <!--        severity="success"-->
    <!--        @click="emits('show-login')"-->
    <!--      />-->
    <!--      <Button-->
    <!--        class="text-white"-->
    <!--        :label="t('account.signup')"-->
    <!--        severity="info"-->
    <!--        @click="emits('show-signup')"-->
    <!--      />-->
    <!--    </div>-->
    <Button
      class="ml-auto text-white"
      icon="pi pi-bars"
      size="large"
      text
      @click="showSidebar = true"
    />
    <Sidebar v-model:visible="showSidebar" :pt="{ closeButton: {class: 'ml-auto'}}">
      <div class="h-full flex flex-col items-start">
        <div class="flex flex-col gap-2">
          <h3 class="underline font-bold">
            {{ t('menu.business') }}
          </h3>
          <router-link
            v-for="item in userItems"
            :key="item.label"
            active-class="font-bold"
            class="menu-item"
            :to="item.path"
            @click="showSidebar = false"
          >
            {{ item.label }}
          </router-link>
          <h3 class="mt-4 underline font-bold">
            {{ t('menu.world') }}
          </h3>
          <router-link
            v-for="item in worldItems"
            :key="item.label"
            active-class="font-bold"
            class="menu-item"
            :to="item.path"
            @click="showSidebar = false"
          >
            {{ item.label }}
          </router-link>
        </div>
        <div class="mt-16 flex justify-between items-center w-full">
          <Button
            v-if="user?.nickName"
            class="p-0"
            :label="t('account.logout')"
            severity="secondary"
            text
            @click="emits('sign-out')"
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
import type { User } from '@/api'
import LangSelect from '@/components/Header/LangSelect.vue'
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
