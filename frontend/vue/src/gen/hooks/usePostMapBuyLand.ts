import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { PostMapBuyLandMutationRequest, PostMapBuyLandMutationResponse, PostMapBuyLand500 } from '../types/PostMapBuyLand.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const postMapBuyLandMutationKey = () => [{ url: '/map/buy_land' }] as const

export type PostMapBuyLandMutationKey = ReturnType<typeof postMapBuyLandMutationKey>

/**
 * @summary Buy land in cell
 * {@link /map/buy_land}
 */
async function postMapBuyLand(data: PostMapBuyLandMutationRequest, config: Partial<RequestConfig<PostMapBuyLandMutationRequest>> = {}) {
  const res = await client<PostMapBuyLandMutationResponse, PostMapBuyLand500, PostMapBuyLandMutationRequest>({
    method: 'POST',
    url: '/map/buy_land',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @summary Buy land in cell
 * {@link /map/buy_land}
 */
export function usePostMapBuyLand(
  options: {
    mutation?: MutationObserverOptions<PostMapBuyLandMutationResponse, PostMapBuyLand500, { data: MaybeRef<PostMapBuyLandMutationRequest> }>;
    client?: Partial<RequestConfig<PostMapBuyLandMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postMapBuyLandMutationKey()

  return useMutation<PostMapBuyLandMutationResponse, PostMapBuyLand500, { data: PostMapBuyLandMutationRequest }>({
    mutationFn: async ({ data }) => {
      return postMapBuyLand(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}