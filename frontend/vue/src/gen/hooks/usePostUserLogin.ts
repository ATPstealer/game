import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { PostUserLoginMutationRequest, PostUserLoginMutationResponse } from '../types/PostUserLogin.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const postUserLoginMutationKey = () => [{ url: '/user/login' }] as const

export type PostUserLoginMutationKey = ReturnType<typeof postUserLoginMutationKey>

/**
 * @description Validate user credentials and provide a JWT token
 * @summary Authenticate a user
 * {@link /user/login}
 */
async function postUserLogin(data: PostUserLoginMutationRequest, config: Partial<RequestConfig<PostUserLoginMutationRequest>> = {}) {
  const res = await client<PostUserLoginMutationResponse, Error, PostUserLoginMutationRequest>({
    method: 'POST',
    url: '/user/login',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @description Validate user credentials and provide a JWT token
 * @summary Authenticate a user
 * {@link /user/login}
 */
export function usePostUserLogin(
  options: {
    mutation?: MutationObserverOptions<PostUserLoginMutationResponse, Error, { data: MaybeRef<PostUserLoginMutationRequest> }>;
    client?: Partial<RequestConfig<PostUserLoginMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postUserLoginMutationKey()

  return useMutation<PostUserLoginMutationResponse, Error, { data: PostUserLoginMutationRequest }>({
    mutationFn: async ({ data }) => {
      return postUserLogin(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}