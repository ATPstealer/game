import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { PostUserLoginMutationRequest, PostUserLoginMutationResponse } from '../types/PostUserLogin.ts'

/**
 * @description Validate user credentials and provide a JWT token
 * @summary Authenticate a user
 * {@link /user/login}
 */
export async function postUserLogin(data: PostUserLoginMutationRequest, config: Partial<RequestConfig<PostUserLoginMutationRequest>> = {}) {
  const res = await client<PostUserLoginMutationResponse, Error, PostUserLoginMutationRequest>({ method: 'POST', url: '/user/login', data, ...config })
  
  return res
}