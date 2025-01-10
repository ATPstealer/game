import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { DeleteUserLoginMutationResponse, DeleteUserLogin500 } from '../types/DeleteUserLogin.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const deleteUserLoginMutationKey = () => [{ url: '/user/login' }] as const

export type DeleteUserLoginMutationKey = ReturnType<typeof deleteUserLoginMutationKey>

/**
 * @description Logout a user by deleting their secure token
 * @summary Logout a user
 * {@link /user/login}
 */
async function deleteUserLogin(config: Partial<RequestConfig> = {}) {
  const res = await client<DeleteUserLoginMutationResponse, DeleteUserLogin500, unknown>({
    method: 'DELETE',
    url: '/user/login',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

/**
 * @description Logout a user by deleting their secure token
 * @summary Logout a user
 * {@link /user/login}
 */
export function useDeleteUserLogin(
  options: {
    mutation?: MutationObserverOptions<DeleteUserLoginMutationResponse, DeleteUserLogin500>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? deleteUserLoginMutationKey()

  return useMutation<DeleteUserLoginMutationResponse, DeleteUserLogin500>({
    mutationFn: async () => {
      return deleteUserLogin(config)
    },
    mutationKey,
    ...mutationOptions
  })
}