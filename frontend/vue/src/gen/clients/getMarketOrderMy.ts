import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { GetMarketOrderMyQueryResponse, GetMarketOrderMy401, GetMarketOrderMy500 } from '../types/GetMarketOrderMy.ts'

/**
 * @summary Get my orders
 * {@link /market/order/my}
 */
export async function getMarketOrderMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMarketOrderMyQueryResponse, GetMarketOrderMy401 | GetMarketOrderMy500, unknown>({
    method: 'GET',
    url: '/market/order/my',
    ...config
  })
  
  return res
}