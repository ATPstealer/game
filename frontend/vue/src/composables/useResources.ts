import { ref } from 'vue'
import { useMyFetch } from '@/composables/useMyFetch'
import type { DataMessage } from '@/types'

export const useResources = () => {
  const moveResource = (payload: any) => {
    const { toX, toY, amount, resourceTypeId, fromX, fromY } = payload
    const dataMessage = ref<DataMessage | null>(null)
    const { onFetchResponse } = useMyFetch(`/resource/move?resource_type_id=${resourceTypeId}&amount=${amount
    }&from_x=${fromX}&from_y=${fromY}&to_x=${toX}&to_y=${toY}`, {
      afterFetch: ctx => {
        dataMessage.value = ctx.data

        return ctx
      }
    }).json()

    return {
      onFetchResponse,
      dataMessage
    }
  }

  return {
    moveResource
  }
}
