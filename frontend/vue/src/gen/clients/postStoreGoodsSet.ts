import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type {
  PostStoreGoodsSetMutationRequest,
  PostStoreGoodsSetMutationResponse,
  PostStoreGoodsSet401,
  PostStoreGoodsSet500
} from '../types/PostStoreGoodsSet.ts'

/**
 * @summary Set prices for goods in the store
 * {@link /store/goods/set}
 */
export async function postStoreGoodsSet(data: PostStoreGoodsSetMutationRequest, config: Partial<RequestConfig<PostStoreGoodsSetMutationRequest>> = {}) {
  const res = await client<PostStoreGoodsSetMutationResponse, PostStoreGoodsSet401 | PostStoreGoodsSet500, PostStoreGoodsSetMutationRequest>({
    method: 'POST',
    url: '/store/goods/set',
    data,
    ...config
  })
  
  return res
}