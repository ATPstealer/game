import { EventHookOn } from '@vueuse/core'
import { Ref, ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'
import type { Message } from '@/types'
import type { Cell } from '@/types/Map/index.interface'
import { CellOwners } from '@/types/Map/index.interface'

export const useMap = () => {
  const getMap = () => {
    const xArray = ref<number[]>([])
    const yArray = ref<number[]>([])

    const { data, onFetchResponse, isFetching } = useMyFetch('/map/', {
      afterFetch: ctx => {
        if (ctx.data) {
          ctx.data = ctx.data.data
        }

        const maxX = ctx.data?.reduce((max: number, obj: Cell) => obj.x > max ? obj.x : max, ctx.data[0].x)
        const maxY = ctx.data?.reduce((max: number, obj: Cell) => obj.y > max ? obj.y : max, ctx.data[0].y)
        const minX = ctx.data?.reduce((min: number, obj: Cell) => obj.x < min ? obj.x : min, ctx.data[0].x)
        const minY = ctx.data?.reduce((min: number, obj: Cell) => obj.y < min ? obj.y : min, ctx.data[0].y)

        if (maxX && minX) {
          for (let i = minX; i <= maxX; i++) {
            xArray.value.push(i)
          }
        }

        if (maxY && minY) {
          for (let i = maxY; i >= minY; i--) {
            yArray.value.push(i)
          }
        }

        return ctx
      }
    }).json()

    return {
      data,
      xArray,
      yArray,
      onFetchResponse,
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