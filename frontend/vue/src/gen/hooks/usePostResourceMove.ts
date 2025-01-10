import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { PostResourceMoveMutationRequest, PostResourceMoveMutationResponse, PostResourceMove401, PostResourceMove500 } from '../types/PostResourceMove.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const postResourceMoveMutationKey = () => [{ url: '/resource/move' }] as const

export type PostResourceMoveMutationKey = ReturnType<typeof postResourceMoveMutationKey>

/**
 * @summary Initiates a resource movement
 * {@link /resource/move}
 */
async function postResourceMove(data?: PostResourceMoveMutationRequest, config: Partial<RequestConfig<PostResourceMoveMutationRequest>> = {}) {
  const res = await client<PostResourceMoveMutationResponse, PostResourceMove401 | PostResourceMove500, PostResourceMoveMutationRequest>({
    method: 'POST',
    url: '/resource/move',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @summary Initiates a resource movement
 * {@link /resource/move}
 */
export function usePostResourceMove(
  options: {
    mutation?: MutationObserverOptions<
      PostResourceMoveMutationResponse,
      PostResourceMove401 | PostResourceMove500,
      { data?: MaybeRef<PostResourceMoveMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostResourceMoveMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postResourceMoveMutationKey()

  return useMutation<PostResourceMoveMutationResponse, PostResourceMove401 | PostResourceMove500, { data?: PostResourceMoveMutationRequest }>({
    mutationFn: async ({ data }) => {
      return postResourceMove(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}