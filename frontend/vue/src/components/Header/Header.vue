<template>
  <header>
    <div class="bg-gray-500 h-8 py-8 px-4 flex">
      <HeaderDesktop
        :user="user"
        :user-items="userItems"
        :world-items="worldItems"
        @show-login="loginModal = true"
        @show-signup="signUpModal = true"
        @sign-out="signOut"
        class="hidden md:flex"
      />
      <HeaderMobile
        :user="user"
        :user-items="userItems"
        :world-items="worldItems"
        @show-login="loginModal = true"
        @show-signup="signUpModal = true"
        @sign-out="signOut"
        class="md:hidden"
      />
    </div>
    <Dialog
      v-model:visible="loginModal"
      modal
      :header="t('account.login')"
      :style="{ width: '25rem' }"
      :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    >
      <LoginForm @close="close" @sign-up="loginModal = false; signUpModal = true" />
    </Dialog>
    <Dialog
      v-model:visible="signUpModal"
      modal
      :header="t('account.signup')"
      :style="{ width: '25rem' }"
      :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    >
      <SignUpForm @close="close" @log-in="signUpModal = false; loginModal = true" />
    </Dialog>
  </header>
</template>

<script setup lang="ts">
import Dialog from 'primevue/dialog'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import LoginForm from '@/components/Common/LoginForm.vue'
import SignUpForm from '@/components/Common/SignUpForm.vue'
import { userPages, worldPages } from '@/components/Header/constants'
import HeaderDesktop from '@/components/Header/HeaderDesktop.vue'
import HeaderMobile from '@/components/Header/HeaderMobile.vue'
import { useGetData } from '@/composables/useGetData'
import { useUser } from '@/composables/useUser'
import type { User } from '@/types'
import type { MenuItem } from '@/types/Header/index.interface'

const { t } = useI18n()

const user = ref<User>()
const loginModal = ref<boolean>(false)
const signUpModal = ref<boolean>(false)

const userItems = computed<MenuItem[]>(() =>
  userPages.map(page => {
    return {
      label: t(`${page}.title`),
      path: page === 'main'? '/' : `/${page}`
    }
  })
)

const worldItems = computed<MenuItem[]>(() =>
  worldPages.map(page => {
    return {
      label: t(`${page}.title`),
      path: `/${page}`
    }
  })
)

const router = useRouter()
const { logOut } = useUser()

const { data, onFetchResponse } = useGetData<User>('/user/data')
onFetchResponse(() => {
  user.value = data.value
})

const close = () => {
  loginModal.value = false
  router.go(0)
}

const signOut = () => {
  logOut()
  router.go(0)
}

</script>

<style scoped>

</style>