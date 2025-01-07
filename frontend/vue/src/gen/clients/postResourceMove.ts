import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { PostResourceMoveMutationRequest, PostResourceMoveMutationResponse, PostResourceMove401, PostResourceMove500 } from '../types/PostResourceMove.ts'

/**
 * @summary Initiates a resource movement
 * {@link /resource/move}
 */
export async function postResourceMove(data?: PostResourceMoveMutationRequest, config: Partial<RequestConfig<PostResourceMoveMutationRequest>> = {}) {
  const res = await client<PostResourceMoveMutationResponse, PostResourceMove401 | PostResourceMove500, PostResourceMoveMutationRequest>({
    method: 'POST',
    url: '/resource/move',
    data,
    ...config
  })

  return res.data
}