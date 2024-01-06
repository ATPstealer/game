import type { EventHookOn } from '@vueuse/core'
import { type Ref, ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'
import type { Message } from '@/types'
import type { Cell, CellOwners } from '@/types/Map/index.interface'

export const useMap = () => {
  const getMap = (): {data: Ref<Cell[]>; isFetching: Ref<boolean>} => {
    const { data, isFetching } = useMyFetch('/map/', {
      afterFetch: ctx => {
        if (ctx.data) {
          ctx.data = ctx.data.data
        }

        return ctx
      }
    }).json()

    return {
      data,
      isFetching
    }
  }

  const getCellOwners = ({ x, y }: {x: number; y: number}): {data: Ref<CellOwners[]>; onFetchResponse: EventHookOn<Response>; isFetching: Ref<boolean> } => {
    const { data, onFetchResponse, isFetching } = useMyFetch<CellOwners[]>(`/map/cell_owners?x=${x}&y=${y}`, {
      afterFetch: ctx => {
        ctx.data = ctx.data.data

        return ctx
      }
    }).json()

    return {
      data,
      onFetchResponse,
      isFetching
    }
  }

  const buyCellSquare = (payload: {x: number; y: number; square: number}) => {
    const dataMessage = ref<Message | null>(null)
    const { data, onFetchResponse } = useMyFetch('/map/buy_land', {
      afterFetch: ctx => {
        dataMessage.value = {
          text: ctx.data.text,
          status: ctx.data.status
        }

        return ctx
      }
    }).post(payload).json()

    return {
      data,
      dataMessage,
      onFetchResponse
    }
  }

  return {
    getMap,
    getCellOwners,
    buyCellSquare
  }
}