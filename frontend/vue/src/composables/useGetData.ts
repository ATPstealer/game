import type { EventHookOn } from '@vueuse/core'
import type { Ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'

export const useGetData = <T>(path: string): {data: Ref<T>; onFetchResponse: EventHookOn<Response>; isFetching: Ref<boolean>; execute: () => void} => {
  const { data, onFetchResponse, isFetching, execute } = useMyFetch(path,
    {
      afterFetch: ctx => {
        ctx.data = ctx.data.data as T

        return ctx
      }
    }).json()

  return {
    data,
    onFetchResponse,
    isFetching,
    execute
  }
}
