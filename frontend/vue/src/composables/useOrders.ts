import { useFetch } from '@vueuse/core'
import type { EventHookOn } from '@vueuse/core'
import { computed } from 'vue'
import type { Ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'
import type { BackData, Order } from '@/types'
import type { MarketParams } from '@/types/Resources/index.interface'

export const useOrders = () => {
  const closeOrder = (orderId: string) => {
    const { data, onFetchResponse, isFetching } = useFetch(`${import.meta.env.VITE_API}/market/order/close?order_id=${orderId}`,
      { credentials: 'include' }).delete()

    return {
      data,
      onFetchResponse,
      isFetching
    }
  }

  const getOrders = (params: MarketParams): {data: Ref<Order[]>; onFetchResponse: EventHookOn<Response>; execute: () => void} => {
    const url = computed(() => {
      // eslint-disable-next-line
      // @ts-expect-error
      const p = new URLSearchParams(params).toString()

      return `/market/order/get?${p}&orderField=priceForUnit&order=-1`
    })

    const { data, onFetchResponse, execute } = useMyFetch(url,
      {
        afterFetch: ctx => {
          if (ctx.data.data) {
            ctx.data = ctx.data.data
          }
          else ctx.data = []

          return ctx
        },
        refetch: true
      }
    ).json()

    return {
      data,
      onFetchResponse,
      execute
    }
  }

  const executeOrder = (payload: {orderID: string; amount: number}): {data: Ref<BackData>; onFetchResponse: EventHookOn<Response>} => {
    const { data, onFetchResponse } = useMyFetch('/market/order/execute').post(payload).json()

    return {
      data,
      onFetchResponse
    }
  }

  const createOrder = (payload: any): {data: Ref<BackData>; onFetchResponse: EventHookOn<Response>} => {
    const { data, onFetchResponse } = useMyFetch('/market/order/create').post(payload).json()

    return {
      data,
      onFetchResponse
    }
  }

  return {
    closeOrder,
    getOrders,
    executeOrder,
    createOrder
  }
}
