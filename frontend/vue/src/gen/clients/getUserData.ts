import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetUserDataQueryResponse, GetUserData401 } from '../types/GetUserData.ts'

/**
 * @summary Get user data
 * {@link /user/data}
 */
export async function getUserData(config: Partial<RequestConfig> = {}) {
  const res = await client<GetUserDataQueryResponse, GetUserData401, unknown>({ method: 'GET', url: '/user/data', ...config })

  return res.data
}