import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostStoreGoodsSetMutationRequest,
  PostStoreGoodsSetMutationResponse,
  PostStoreGoodsSet401,
  PostStoreGoodsSet500
} from '../types/PostStoreGoodsSet.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const postStoreGoodsSetMutationKey = () => [{ url: '/store/goods/set' }] as const

export type PostStoreGoodsSetMutationKey = ReturnType<typeof postStoreGoodsSetMutationKey>

/**
 * @summary Set prices for goods in the store
 * {@link /store/goods/set}
 */
async function postStoreGoodsSet(data: PostStoreGoodsSetMutationRequest, config: Partial<RequestConfig<PostStoreGoodsSetMutationRequest>> = {}) {
  const res = await client<
    PostStoreGoodsSetMutationResponse,
    ResponseErrorConfig<PostStoreGoodsSet401 | PostStoreGoodsSet500>,
    PostStoreGoodsSetMutationRequest
  >({ method: 'POST', url: '/store/goods/set', data, ...config })
  
  return res.data
}

/**
 * @summary Set prices for goods in the store
 * {@link /store/goods/set}
 */
export function usePostStoreGoodsSet(
  options: {
    mutation?: MutationObserverOptions<
      PostStoreGoodsSetMutationResponse,
      ResponseErrorConfig<PostStoreGoodsSet401 | PostStoreGoodsSet500>,
      { data: MaybeRef<PostStoreGoodsSetMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostStoreGoodsSetMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postStoreGoodsSetMutationKey()

  return useMutation<
    PostStoreGoodsSetMutationResponse,
    ResponseErrorConfig<PostStoreGoodsSet401 | PostStoreGoodsSet500>,
    { data: PostStoreGoodsSetMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postStoreGoodsSet(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}