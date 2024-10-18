<template>
  <header>
    <div class="bg-gray-500 h-8 py-8 px-4 flex">
      <HeaderDesktop
        class="hidden md:flex"
        :help-items="helpItems"
        :user="user"
        :user-items="userItems"
        :world-items="worldItems"
        @show-login="loginModal = true"
        @show-signup="signUpModal = true"
        @sign-out="signOut"
      />
      <HeaderMobile
        class="md:hidden"
        :help-items="helpItems"
        :user="user"
        :user-items="userItems"
        :world-items="worldItems"
        @show-login="loginModal = true"
        @show-signup="signUpModal = true"
        @sign-out="signOut"
      />
    </div>
    <Dialog
      v-model:visible="loginModal"
      :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
      :header="t('account.login')"
      modal
      :style="{ width: '25rem' }"
    >
      <LoginForm @close="close" @sign-up="loginModal = false; signUpModal = true" />
    </Dialog>
    <Dialog
      v-model:visible="signUpModal"
      :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
      :header="t('account.signup')"
      modal
      :style="{ width: '25rem' }"
    >
      <SignUpForm @close="close" @log-in="signUpModal = false; loginModal = true" />
    </Dialog>
  </header>
</template>

<script setup lang="ts">
import { useMutation, useQuery } from '@tanstack/vue-query'
import { storeToRefs } from 'pinia'
import Dialog from 'primevue/dialog'
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import type { User } from '@/api'
import { deleteUserLoginMutation, getUserDataOptions } from '@/api/@tanstack/vue-query.gen'
import LoginForm from '@/components/Common/LoginForm.vue'
import SignUpForm from '@/components/Common/SignUpForm.vue'
import { helpPages, userPages, worldPages } from '@/components/Header/constants'
import HeaderDesktop from '@/components/Header/HeaderDesktop.vue'
import HeaderMobile from '@/components/Header/HeaderMobile.vue'
import { useUserStore } from '@/stores/userStore'
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

const helpItems = computed<MenuItem[]>(() =>
  helpPages.map(page => {
    return {
      label: t(`${page}.title`),
      path: `/${page}`
    }
  })
)

const router = useRouter()

const { data: userFetchedData, isSuccess, refetch } = useQuery({
  ...getUserDataOptions(),
  enabled: false,
  select: data => data.data
})

const { userData } = storeToRefs(useUserStore())
const { setUserData } = useUserStore()

if (!userData.value) {
  refetch()
}

const logout = useMutation({
  ...deleteUserLoginMutation(),
  onSuccess: () => {
    localStorage.setItem('invalid', 'true')
    router.go(0)
  }
})

const close = () => {
  loginModal.value = false
  router.go(0)
}

const signOut = () => {
  logout.mutate({})
}

watch(isSuccess, () => {
  if (isSuccess.value) {
    user.value = userFetchedData!.value
    setUserData(userFetchedData!.value)
  }
})

</script>

<style scoped>

</style>