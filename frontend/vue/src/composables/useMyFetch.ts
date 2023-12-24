import { createFetch } from '@vueuse/core'

export const useMyFetch = createFetch({
  baseUrl: `${import.meta.env.VITE_API}`,
  fetchOptions: {
    credentials: 'include'
  }
})