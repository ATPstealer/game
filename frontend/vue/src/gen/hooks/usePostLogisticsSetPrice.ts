import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostLogisticsSetPriceMutationRequest,
  PostLogisticsSetPriceMutationResponse,
  PostLogisticsSetPrice401,
  PostLogisticsSetPrice500
} from '../types/PostLogisticsSetPrice.ts'

export const postLogisticsSetPriceMutationKey = () => [{ url: '/logistics/set_price' }] as const

export type PostLogisticsSetPriceMutationKey = ReturnType<typeof postLogisticsSetPriceMutationKey>

/**
 * @summary Set the logistics price
 * {@link /logistics/set_price}
 */
async function postLogisticsSetPrice(data?: PostLogisticsSetPriceMutationRequest, config: Partial<RequestConfig<PostLogisticsSetPriceMutationRequest>> = {}) {
  const res = await client<PostLogisticsSetPriceMutationResponse, PostLogisticsSetPrice401 | PostLogisticsSetPrice500, PostLogisticsSetPriceMutationRequest>({
    method: 'POST',
    url: '/logistics/set_price',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @summary Set the logistics price
 * {@link /logistics/set_price}
 */
export function usePostLogisticsSetPrice(
  options: {
    mutation?: MutationObserverOptions<
      PostLogisticsSetPriceMutationResponse,
      PostLogisticsSetPrice401 | PostLogisticsSetPrice500,
      { data?: MaybeRef<PostLogisticsSetPriceMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostLogisticsSetPriceMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postLogisticsSetPriceMutationKey()

  return useMutation<
    PostLogisticsSetPriceMutationResponse,
    PostLogisticsSetPrice401 | PostLogisticsSetPrice500,
    { data?: PostLogisticsSetPriceMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postLogisticsSetPrice(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}