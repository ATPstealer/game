import { createFetch } from '@vueuse/core'
import router from '@/router'

export const useMyFetch = createFetch({
  baseUrl: `${import.meta.env.VITE_API}`,
  options: {
    onFetchError: async ctx => {
      await router.isReady()

      const currentRoute = router.currentRoute.value.fullPath

      if (ctx.response?.status === 401 && currentRoute !== '/login') {
        localStorage.setItem('prev', currentRoute)
        await router.push({ name: 'Login' })
      }

      return ctx
    }
  },
  fetchOptions: {
    credentials: 'include'
  }
})