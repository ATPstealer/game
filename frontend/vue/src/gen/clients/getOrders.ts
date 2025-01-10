import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { GetOrdersQueryResponse, GetOrdersQueryParams, GetOrders500 } from '../types/GetOrders.ts'

/**
 * @summary Fetches orders based on various query parameters
 * {@link /orders}
 */
export async function getOrders(params?: GetOrdersQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetOrdersQueryResponse, GetOrders500, unknown>({ method: 'GET', url: '/orders', params, ...config })
  
  return res
}