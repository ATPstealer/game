import type { EventHookOn } from '@vueuse/core'
import type { Ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useMyFetch } from '@/composables/useMyFetch'
import { useNotify } from '@/composables/useNotify'
import type { BackData } from '@/types'
import type { ConstructBuildingPayload, SearchBuildingParams } from '@/types/Buildings/index.interface'

export const useBuildings = () => {
  const  { setWarning } = useNotify()
  const { t } = useI18n()

  const constructBuilding = (payload: ConstructBuildingPayload): {data: Ref<BackData>; onFetchFinally: EventHookOn<Response>} => {
    const { data, onFetchFinally } = useMyFetch('/building/construct').post(payload).json()

    return {
      data,
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

  const startProduction = (payload: { duration: number; blueprintId: number; buildingId: string }): {data: Ref<BackData>; onFetchResponse: EventHookOn<Response>} => {
    const { data, onFetchResponse } = useMyFetch('/building/start_work').post(payload).json()

    return {
      data,
      onFetchResponse
    }
  }

  const stopWork = (payload: { buildingId: string }): {data: Ref<BackData>; onFetchResponse: EventHookOn<Response>} => {
    const { data, onFetchResponse } = useMyFetch('/building/stop_work').post(payload).json()

    return {
      data,
      onFetchResponse
    }
  }

  const setPrice = (payload: { price: any; resourceTypeId: any; buildingId: string }) => {
    const { onFetchResponse, isFetching } = useMyFetch('/store/goods/set').post(payload).json()

    return {
      onFetchResponse,
      isFetching
    }
  }

  const setHiring = (payload: {buildingId: string; salary: number; hiringNeeds: number}): {data: Ref<BackData>; onFetchResponse: EventHookOn<Response>} => {
    const { data, onFetchResponse } = useMyFetch('/building/hiring').post(payload).json()

    return {
      data,
      onFetchResponse
    }
  }
  const destroyBuilding = (id: string) => {
    const { data, onFetchResponse } = useMyFetch(`/building/destroy?_id=${id}`).delete().json()

    onFetchResponse(() => {
      setWarning(t(`codes.${data.value.code.toString()}`))
    })

    return {
      onFetchResponse
    }
  }

  return {
    constructBuilding,
    getBuildings,
    startProduction,
    stopWork,
    setPrice,
    setHiring,
    destroyBuilding
  }
}
