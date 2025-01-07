import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetUsernamesByPrefixQueryResponse, GetUsernamesByPrefixQueryParams } from '../types/GetUsernamesByPrefix.ts'

/**
 * @description Retrieve a list of usernames that match the given prefix
 * @summary Get usernames by prefix
 * {@link /data/users_by_prefix}
 */
export async function getUsernamesByPrefix(params?: GetUsernamesByPrefixQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetUsernamesByPrefixQueryResponse, Error, unknown>({ method: 'GET', url: '/data/users_by_prefix', params, ...config })
  
  return res.data
}