<template>
  <div class="flex justify-between items-center w-full">
    <div class="flex gap-8">
      <Button
        text
        @click="toggleUserItems"
        :label="t('menu.business')"
        aria-haspopup="true"
        aria-controls="overlay_user"
        class="header-item"
      />
      <Menu
        ref="userMenu"
        id="overlay_user"
        :model="userItems"
        :popup="true"
      >
        <template #item="{item}">
          <router-link
            active-class="font-bold"
            class="menu-item"
            :to="item.path"
          >
            {{ item.label }}
          </router-link>
        </template>
      </Menu>
      <Button
        text
        @click="toggleWorldItems"
        :label="t('menu.world')"
        aria-haspopup="true"
        aria-controls="overlay_world"
        class="header-item"
      />
      <Menu
        ref="worldMenu"
        id="overlay_world"
        :model="worldItems"
        :popup="true"
      >
        <template #item="{item}">
          <router-link
            active-class="font-bold"
            class="menu-item"
            :to="item.path"
          >
            {{ item.label }}
          </router-link>
        </template>
      </Menu>
    </div>
    <div class="flex items-center gap-8">
      <div v-if="!user?.nickName" class="flex gap-4">
        <Button
          label="Login"
          @click="emits('show-login')"
          class="text-white"
          severity="success"
        />
        <Button
          label="Sign Up"
          @click="emits('show-signup')"
          class="text-white"
          severity="info"
        />
      </div>
      <div v-else class="flex items-center gap-4">
        <p class="header-item">
          {{ user?.nickName }}
        </p>
        <p class="header-item">
          {{ moneyFormat(user?.money) }}
        </p>
        <Button
          text
          label="Logout"
          class="header-item p-0"
          @click="emits('sign-out')"
        />
      </div>
      <LangSelect />
    </div>
  </div>
</template>

<script setup lang="ts">

import Button from 'primevue/button'
import Menu from 'primevue/menu'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import LangSelect from '@/components/Header/LangSelect.vue'
import { User } from '@/types'
import { MenuItem } from '@/types/Header/index.interface'
import {moneyFormat} from "@/utils/moneyFormat";

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

const userMenu = ref()
const worldMenu = ref()

const { t } = useI18n()

const toggleUserItems = (event) => {
  userMenu.value.toggle(event)
}
const toggleWorldItems = (event) => {
  worldMenu.value.toggle(event)
}
</script>

<style scoped>
.header-item {
  @apply text-white font-bold;
}

.menu-item {
  @apply block pl-4 py-1 no-underline text-gray-800;
}
</style>
