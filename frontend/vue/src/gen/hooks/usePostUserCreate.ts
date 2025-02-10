import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { PostUserCreateMutationRequest, PostUserCreateMutationResponse, PostUserCreate500 } from '../types/PostUserCreate.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const postUserCreateMutationKey = () => [{ url: '/user/create' }] as const

export type PostUserCreateMutationKey = ReturnType<typeof postUserCreateMutationKey>

/**
 * @summary Create a new user
 * {@link /user/create}
 */
async function postUserCreate(data: PostUserCreateMutationRequest, config: Partial<RequestConfig<PostUserCreateMutationRequest>> = {}) {
  const res = await client<PostUserCreateMutationResponse, ResponseErrorConfig<PostUserCreate500>, PostUserCreateMutationRequest>({
    method: 'POST',
    url: '/user/create',
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
    mutation?: MutationObserverOptions<
      PostUserCreateMutationResponse,
      ResponseErrorConfig<PostUserCreate500>,
      { data: MaybeRef<PostUserCreateMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostUserCreateMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postUserCreateMutationKey()

  return useMutation<PostUserCreateMutationResponse, ResponseErrorConfig<PostUserCreate500>, { data: PostUserCreateMutationRequest }>({
    mutationFn: async ({ data }) => {
      return postUserCreate(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}