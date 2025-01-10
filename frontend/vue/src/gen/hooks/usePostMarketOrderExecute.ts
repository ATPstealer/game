import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostMarketOrderExecuteMutationRequest,
  PostMarketOrderExecuteMutationResponse,
  PostMarketOrderExecute401,
  PostMarketOrderExecute500
} from '../types/PostMarketOrderExecute.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const postMarketOrderExecuteMutationKey = () => [{ url: '/market/order/execute' }] as const

export type PostMarketOrderExecuteMutationKey = ReturnType<typeof postMarketOrderExecuteMutationKey>

/**
 * @summary Partially execute an  order
 * {@link /market/order/execute}
 */
async function postMarketOrderExecute(
  data?: PostMarketOrderExecuteMutationRequest,
  config: Partial<RequestConfig<PostMarketOrderExecuteMutationRequest>> = {}
) {
  const res = await client<
    PostMarketOrderExecuteMutationResponse,
    PostMarketOrderExecute401 | PostMarketOrderExecute500,
    PostMarketOrderExecuteMutationRequest
  >({ method: 'POST', url: '/market/order/execute', baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2', data, ...config })
  
  return res.data
}

/**
 * @summary Partially execute an  order
 * {@link /market/order/execute}
 */
export function usePostMarketOrderExecute(
  options: {
    mutation?: MutationObserverOptions<
      PostMarketOrderExecuteMutationResponse,
      PostMarketOrderExecute401 | PostMarketOrderExecute500,
      { data?: MaybeRef<PostMarketOrderExecuteMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostMarketOrderExecuteMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postMarketOrderExecuteMutationKey()

  return useMutation<
    PostMarketOrderExecuteMutationResponse,
    PostMarketOrderExecute401 | PostMarketOrderExecute500,
    { data?: PostMarketOrderExecuteMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postMarketOrderExecute(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}