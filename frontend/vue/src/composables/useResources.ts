import type { EventHookOn } from '@vueuse/core'
import type { Ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'
import type { BackData } from '@/types'
import type { ResourceMovePayload } from '@/types/Resources/index.interface'

export const useResources = () => {
  const moveResource = (payload: ResourceMovePayload): {data: Ref<BackData>; onFetchResponse: EventHookOn<Response>} => {
    const { data, onFetchResponse } = useMyFetch('/resource/move').post(payload).json()

    return {
      data,
      onFetchResponse
    }
  }

  return {
    moveResource
  }
}
