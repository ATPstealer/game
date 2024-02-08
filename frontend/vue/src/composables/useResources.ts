import { ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'
import type { DataMessage } from '@/types'
import type { ResourceMovePayload } from '@/types/Resources/index.interface'

export const useResources = () => {
  const moveResource = (payload: ResourceMovePayload) => {
    const dataMessage = ref<DataMessage | null>(null)
    const { onFetchResponse } = useMyFetch('/resource/move', {
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
    moveResource
  }
}
