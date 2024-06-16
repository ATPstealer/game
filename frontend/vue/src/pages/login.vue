<template>
  <Card class="max-w-lg m-auto">
    <template #content>
      <LoginForm
        v-if="isLogin"
        @close="redirect"
        @sign-up="isLogin = false"
      />
      <SignUpForm
        v-else
        @close="redirect"
        @log-in="isLogin = true"
      />
    </template>
  </Card>
</template>

<script setup lang="ts">
import Card from 'primevue/card'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import LoginForm from '@/components/Common/LoginForm.vue'
import SignUpForm from '@/components/Common/SignUpForm.vue'

const router = useRouter()

const isLogin = ref<boolean>(true)

const redirect = async () => {
  const prev = localStorage.getItem('prev')
  await router.push(prev)
  router.go(0)
}
</script>

<style scoped>

</style>