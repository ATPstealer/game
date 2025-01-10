import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { PostMapBuyLandMutationRequest, PostMapBuyLandMutationResponse, PostMapBuyLand500 } from '../types/PostMapBuyLand.ts'

/**
 * @summary Buy land in cell
 * {@link /map/buy_land}
 */
export async function postMapBuyLand(data: PostMapBuyLandMutationRequest, config: Partial<RequestConfig<PostMapBuyLandMutationRequest>> = {}) {
  const res = await client<PostMapBuyLandMutationResponse, PostMapBuyLand500, PostMapBuyLandMutationRequest>({
    method: 'POST',
    url: '/map/buy_land',
    data,
    ...config
  })
  
  return res
}