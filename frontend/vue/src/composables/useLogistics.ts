import { useGetData } from '@/composables/useGetData'
import { useMyFetch } from '@/composables/useMyFetch'
import type { LogisticHub } from '@/types/Resources/index.interface'

export const useLogistics = () => {
  const getHubs = () => {
    const { data, onFetchResponse } = useGetData<LogisticHub[]>('/resource/logistics')

    return {
      data,
      onFetchResponse
    }
  }

  const setHubPrice = ({ buildingId, price }: { buildingId: string; price: number }) => {
    const { data, onFetchResponse } = useMyFetch('/logistics/set_price').post({ price, buildingId }).json()

    return {
      data,
      onFetchResponse
    }
  }

  return {
    getHubs,
    setHubPrice
  }
}