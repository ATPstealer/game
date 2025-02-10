import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostMarketOrderCreateMutationRequest,
  PostMarketOrderCreateMutationResponse,
  PostMarketOrderCreate401,
  PostMarketOrderCreate500
} from '../types/PostMarketOrderCreate.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const postMarketOrderCreateMutationKey = () => [{ url: '/market/order/create' }] as const

export type PostMarketOrderCreateMutationKey = ReturnType<typeof postMarketOrderCreateMutationKey>

/**
 * @summary Create a new market order
 * {@link /market/order/create}
 */
async function postMarketOrderCreate(data: PostMarketOrderCreateMutationRequest, config: Partial<RequestConfig<PostMarketOrderCreateMutationRequest>> = {}) {
  const res = await client<
    PostMarketOrderCreateMutationResponse,
    ResponseErrorConfig<PostMarketOrderCreate401 | PostMarketOrderCreate500>,
    PostMarketOrderCreateMutationRequest
  >({ method: 'POST', url: '/market/order/create', data, ...config })
  
  return res.data
}

/**
 * @summary Create a new market order
 * {@link /market/order/create}
 */
export function usePostMarketOrderCreate(
  options: {
    mutation?: MutationObserverOptions<
      PostMarketOrderCreateMutationResponse,
      ResponseErrorConfig<PostMarketOrderCreate401 | PostMarketOrderCreate500>,
      { data: MaybeRef<PostMarketOrderCreateMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostMarketOrderCreateMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postMarketOrderCreateMutationKey()

  return useMutation<
    PostMarketOrderCreateMutationResponse,
    ResponseErrorConfig<PostMarketOrderCreate401 | PostMarketOrderCreate500>,
    { data: PostMarketOrderCreateMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postMarketOrderCreate(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}