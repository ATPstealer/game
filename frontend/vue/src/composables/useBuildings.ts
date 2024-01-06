import { ref } from 'vue'
import type { Ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'
import type { Message } from '@/types'
import type { ConstructBuildingPayload, SearchBuildingParams } from '@/types/Buildings/index.interface'

export const useBuildings = () => {
  const constructBuilding = (payload: ConstructBuildingPayload) => {
    const dataMessage = ref<Message | null>(null)
    const { onFetchFinally } = useMyFetch('/building/construct',
      {
        afterFetch: ctx => {
          dataMessage.value = {
            text: ctx.data.text,
            status: ctx.data.status
          }

          return ctx
        }
      }
    ).post(payload).json()

    return {
      dataMessage,
      onFetchFinally
    }
  }

  const getBuildings = (searchParams: Ref<SearchBuildingParams>) => {
    const { data, isFetching } = useMyFetch('/building/get', {
      afterFetch: ctx => {
        if (ctx.data.data) {
          ctx.data = ctx.data.data
        } else {
          ctx.data = []
        }

        return ctx
      },
      refetch: true
    }).post(searchParams).json()

    return {
      data,
      isFetching
    }
  }

  const startProduction = (payload: {buildingId: number; blueprintId: number; duration: number}) => {
    const dataMessage = ref<Message | null>(null)
    const { onFetchResponse } = useMyFetch('/building/start_work', {
      afterFetch: ctx => {
        dataMessage.value = {
          text: ctx.data.text,
          status: ctx.data.status
        }

        return ctx
      }
    }).post(payload).json()

    return {
      dataMessage,
      onFetchResponse
    }
  }

  const setPrice = (payload: {buildingId: number; resourceTypeId: number; price: number}) => {
    const dataMessage = ref<Message | null>(null)

    const { onFetchResponse, isFetching } = useMyFetch(`/store/goods/set?building_id=${payload.buildingId}&resource_type_id=${payload.resourceTypeId}&price=${payload.price}`, {
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
      onFetchResponse,
      isFetching
    }
  }

  return {
    constructBuilding,
    getBuildings,
    startProduction,
    setPrice
  }
}
