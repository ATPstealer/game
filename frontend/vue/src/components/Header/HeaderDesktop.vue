<template>
  <div class="flex justify-between items-center w-full">
    <div class="flex gap-8">
      <Button
        aria-controls="overlay_user"
        aria-haspopup="true"
        class="header-item"
        :label="t('menu.business')"
        text
        @click="toggleUserItems"
      />
      <Menu
        id="overlay_user"
        ref="userMenu"
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
        aria-controls="overlay_world"
        aria-haspopup="true"
        class="header-item"
        :label="t('menu.world')"
        text
        @click="toggleWorldItems"
      />
      <Menu
        id="overlay_world"
        ref="worldMenu"
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
      <Button
        aria-controls="overlay_help"
        aria-haspopup="true"
        class="header-item"
        :label="t('menu.help')"
        text
        @click="toggleHelpItems"
      />
      <Menu
        id="overlay_help"
        ref="helpMenu"
        :model="helpItems"
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
      <!--      <div v-if="!user?.nickName" class="flex gap-4">-->
      <!--        <Button-->
      <!--          class="text-white"-->
      <!--          :label="t('account.login')"-->
      <!--          severity="success"-->
      <!--          @click="emits('show-login')"-->
      <!--        />-->
      <!--        <Button-->
      <!--          class="text-white"-->
      <!--          :label="t('account.signup')"-->
      <!--          severity="info"-->
      <!--          @click="emits('show-signup')"-->
      <!--        />-->
      <!--      </div>-->
      <div class="flex items-center gap-4">
        <p v-if="user?.nickName" class="header-item">
          {{ user?.nickName }}
        </p>
        <p v-if="user?.money" class="header-item">
          {{ moneyFormat(user?.money) }}
        </p>
        <Button
          v-if="user?.nickName"
          class="header-item p-0"
          :label="t('account.logout')"
          text
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
import type { User } from '@/api'
import LangSelect from '@/components/Header/LangSelect.vue'
import type { MenuItem } from '@/types/Header/index.interface'
import { moneyFormat } from '@/utils/moneyFormat'

interface Props {
  user: User | undefined;
  userItems: MenuItem[];
  worldItems: MenuItem[];
  helpItems: MenuItem[];
}

defineProps<Props>()

const emits = defineEmits<{
  (e: 'show-login'): void;
  (e: 'show-signup'): void;
  (e: 'sign-out'): void;
}>()

const userMenu = ref()
const worldMenu = ref()
const helpMenu = ref()

const { t } = useI18n()

const toggleUserItems = (event: any) => {
  userMenu.value.toggle(event)
}
const toggleWorldItems = (event: any) => {
  worldMenu.value.toggle(event)
}

const toggleHelpItems = (event: any) => {
  helpMenu.value.toggle(event)
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
