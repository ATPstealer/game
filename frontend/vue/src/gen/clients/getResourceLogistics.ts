import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetResourceLogisticsQueryResponse, GetResourceLogisticsQueryParams, GetResourceLogistics500 } from '../types/GetResourceLogistics.ts'

/**
 * @summary Get the logistics capacity in cell
 * {@link /resource/logistics}
 */
export async function getResourceLogistics(params?: GetResourceLogisticsQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetResourceLogisticsQueryResponse, GetResourceLogistics500, unknown>({
    method: 'GET',
    url: '/resource/logistics',
    params,
    ...config
  })

  return res.data
}