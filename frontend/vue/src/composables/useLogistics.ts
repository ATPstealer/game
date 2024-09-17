import { useGetData } from '@/composables/useGetData'
import type { LogisticHub } from '@/types/Resources/index.interface'

export const useLogistics = () => {
  const getHubs = () => {
    const { data, onFetchResponse } = useGetData<LogisticHub[]>('/resource/logistics')

    return {
      data,
      onFetchResponse
    }
  }

  return {
    getHubs
  }
}