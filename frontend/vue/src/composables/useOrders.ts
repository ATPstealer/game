import { useFetch } from '@vueuse/core'
import type { EventHookOn } from '@vueuse/core'
import { computed, ref } from 'vue'
import type { Ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'
import type { DataMessage, Order } from '@/types'
import type { MarketParams } from '@/types/Resources/index.interface'

export const useOrders = () => {
  const closeOrder = (orderId: number) => {
    const { data, onFetchResponse, isFetching } = useFetch(`${import.meta.env.VITE_API}/market/order/close?order_id=${orderId}`,
      { credentials: 'include' })

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

      return `/market/order/get?${p}&order_field=price_for_unit&order=DESC`
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

  const executeOrder = (id: number) => {
    const dataMessage = ref<DataMessage | null>(null)

    const { onFetchResponse } = useMyFetch(`/market/order/execute?order_id=${id}`, {
      afterFetch: ctx => {
        dataMessage.value = {
          text: ctx.data.text,
          status: ctx.data.status
        }

        return ctx
      }
    }).json()

    return {
      dataMessage,
      onFetchResponse
    }
  }

  const createOrder = (payload: any) => {
    const dataMessage = ref<DataMessage | null>(null)
    const { onFetchResponse } = useMyFetch('/market/order/create', {
      afterFetch: ctx => {
        dataMessage.value = ctx.data

        return ctx
      }
    }).post(payload).json()

    return {
      onFetchResponse,
      dataMessage
    }
  }

  return {
    closeOrder,
    getOrders,
    executeOrder,
    createOrder
  }
}
