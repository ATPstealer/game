import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type {
  PostMarketOrderExecuteMutationRequest,
  PostMarketOrderExecuteMutationResponse,
  PostMarketOrderExecute401,
  PostMarketOrderExecute500
} from '../types/PostMarketOrderExecute.ts'

/**
 * @summary Partially execute an  order
 * {@link /market/order/execute}
 */
export async function postMarketOrderExecute(
  data?: PostMarketOrderExecuteMutationRequest,
  config: Partial<RequestConfig<PostMarketOrderExecuteMutationRequest>> = {}
) {
  const res = await client<
    PostMarketOrderExecuteMutationResponse,
    PostMarketOrderExecute401 | PostMarketOrderExecute500,
    PostMarketOrderExecuteMutationRequest
  >({ method: 'POST', url: '/market/order/execute', data, ...config })

  return res.data
}