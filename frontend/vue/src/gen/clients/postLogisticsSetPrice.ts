import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type {
  PostLogisticsSetPriceMutationRequest,
  PostLogisticsSetPriceMutationResponse,
  PostLogisticsSetPrice401,
  PostLogisticsSetPrice500
} from '../types/PostLogisticsSetPrice.ts'

/**
 * @summary Set the logistics price
 * {@link /logistics/set_price}
 */
export async function postLogisticsSetPrice(
  data?: PostLogisticsSetPriceMutationRequest,
  config: Partial<RequestConfig<PostLogisticsSetPriceMutationRequest>> = {}
) {
  const res = await client<PostLogisticsSetPriceMutationResponse, PostLogisticsSetPrice401 | PostLogisticsSetPrice500, PostLogisticsSetPriceMutationRequest>({
    method: 'POST',
    url: '/logistics/set_price',
    data,
    ...config
  })

  return res.data
}