import type { EventHookOn } from '@vueuse/core'
import type { Ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'
import type { BackData } from '@/types'

export const useUser = () => {
  const logIn = (userData: Record<string, string>): {data: Ref<BackData>; onFetchFinally: EventHookOn<Response>} => {
    const { data, onFetchResponse, onFetchFinally }  = useMyFetch('/user/login').post(userData).json()

    onFetchResponse(() => {
      if (data?.value.code <= 0) {
        const ttl = Number(data.value?.data.ttl)
        const date = new Date()
        date.setTime(date.getTime() + ttl*1000)
        document.cookie = `secureToken=${data.value?.data.token};expires=${date};domain=.${import.meta.env.VITE_DOMAIN};path=/`
      }
    })

    return {
      onFetchFinally,
      data
    }
  }

  const logOut = () => {
    useMyFetch('/user/login', {
      onFetchError: ctx => {
        console.log(ctx.error.message)

        return ctx
      }
    }).delete()
  }

  const signUp = (payload: {nickName: string; email: string; password: string}): {data: Ref<BackData>; onFetchFinally: EventHookOn<Response>} => {
    const { data, onFetchFinally } = useMyFetch('/user/create').post(payload).json()

    return {
      data,
      onFetchFinally
    }
  }

  return {
    logIn,
    logOut,
    signUp
  }
}
