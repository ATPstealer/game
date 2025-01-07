import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { DeleteUserLoginMutationResponse, DeleteUserLogin500 } from '../types/DeleteUserLogin.ts'

/**
 * @description Logout a user by deleting their secure token
 * @summary Logout a user
 * {@link /user/login}
 */
export async function deleteUserLogin(config: Partial<RequestConfig> = {}) {
  const res = await client<DeleteUserLoginMutationResponse, DeleteUserLogin500, unknown>({ method: 'DELETE', url: '/user/login', ...config })

  return res.data
}