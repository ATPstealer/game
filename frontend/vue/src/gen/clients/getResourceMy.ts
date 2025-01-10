import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { GetResourceMyQueryResponse, GetResourceMyQueryParams, GetResourceMy401, GetResourceMy500 } from '../types/GetResourceMy.ts'

/**
 * @summary Get user's resources
 * {@link /resource/my}
 */
export async function getResourceMy(params?: GetResourceMyQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetResourceMyQueryResponse, GetResourceMy401 | GetResourceMy500, unknown>({ method: 'GET', url: '/resource/my', params, ...config })
  
  return res
}