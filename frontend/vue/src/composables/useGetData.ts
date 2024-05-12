import type { EventHookOn } from '@vueuse/core'
import type { Ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'

export const useGetData = <T>(path: string, immediate = true): {data: Ref<T>; onFetchResponse: EventHookOn<Response>; isFetching: Ref<boolean>; execute: () => void} => {
  const { data, onFetchResponse, isFetching, execute } = useMyFetch(path,
    {
      // TODO: тупит eslint
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-expect-error
      immediate
    },
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
