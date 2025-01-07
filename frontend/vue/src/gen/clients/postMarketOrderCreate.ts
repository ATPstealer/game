import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type {
  PostMarketOrderCreateMutationRequest,
  PostMarketOrderCreateMutationResponse,
  PostMarketOrderCreate401,
  PostMarketOrderCreate500
} from '../types/PostMarketOrderCreate.ts'

/**
 * @summary Create a new market order
 * {@link /market/order/create}
 */
export async function postMarketOrderCreate(
  data: PostMarketOrderCreateMutationRequest,
  config: Partial<RequestConfig<PostMarketOrderCreateMutationRequest>> = {}
) {
  const res = await client<PostMarketOrderCreateMutationResponse, PostMarketOrderCreate401 | PostMarketOrderCreate500, PostMarketOrderCreateMutationRequest>({
    method: 'POST',
    url: '/market/order/create',
    data,
    ...config
  })

  return res.data
}