import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetResourceMyLogisticsQueryResponse, GetResourceMyLogistics401, GetResourceMyLogistics500 } from '../types/GetResourceMyLogistics.ts'

/**
 * @summary Get user's logistics tasks
 * {@link /resource/my_logistics}
 */
export async function getResourceMyLogistics(config: Partial<RequestConfig> = {}) {
  const res = await client<GetResourceMyLogisticsQueryResponse, GetResourceMyLogistics401 | GetResourceMyLogistics500, unknown>({
    method: 'GET',
    url: '/resource/my_logistics',
    ...config
  })
  
  return res.data
}