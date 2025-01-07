import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { PostUserCreateMutationRequest, PostUserCreateMutationResponse, PostUserCreate500 } from '../types/PostUserCreate.ts'

export const postUserCreateMutationKey = () => [{ url: '/user/create' }] as const

export type PostUserCreateMutationKey = ReturnType<typeof postUserCreateMutationKey>

/**
 * @summary Create a new user
 * {@link /user/create}
 */
async function postUserCreate(data: PostUserCreateMutationRequest, config: Partial<RequestConfig<PostUserCreateMutationRequest>> = {}) {
  const res = await client<PostUserCreateMutationResponse, PostUserCreate500, PostUserCreateMutationRequest>({
    method: 'POST',
    url: '/user/create',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @summary Create a new user
 * {@link /user/create}
 */
export function usePostUserCreate(
  options: {
    mutation?: MutationObserverOptions<PostUserCreateMutationResponse, PostUserCreate500, { data: MaybeRef<PostUserCreateMutationRequest> }>;
    client?: Partial<RequestConfig<PostUserCreateMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postUserCreateMutationKey()

  return useMutation<PostUserCreateMutationResponse, PostUserCreate500, { data: PostUserCreateMutationRequest }>({
    mutationFn: async ({ data }) => {
      return postUserCreate(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}