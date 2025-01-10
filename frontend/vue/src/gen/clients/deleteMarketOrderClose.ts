import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type {
  DeleteMarketOrderCloseMutationResponse,
  DeleteMarketOrderCloseQueryParams,
  DeleteMarketOrderClose401,
  DeleteMarketOrderClose500
} from '../types/DeleteMarketOrderClose.ts'

/**
 * @summary Close user's order
 * {@link /market/order/close}
 */
export async function deleteMarketOrderClose(params: DeleteMarketOrderCloseQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<DeleteMarketOrderCloseMutationResponse, DeleteMarketOrderClose401 | DeleteMarketOrderClose500, unknown>({
    method: 'DELETE',
    url: '/market/order/close',
    params,
    ...config
  })
  
  return res
}