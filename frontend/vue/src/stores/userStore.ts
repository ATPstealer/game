import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/gen'

export const useUserStore = defineStore('userStore', () => {
  const userData = ref<User>()

  const setUserData = (user: User) => {
    userData.value = user
  }

  return {
    userData,
    setUserData
  }
})