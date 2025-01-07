import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { PostUserCreateMutationRequest, PostUserCreateMutationResponse, PostUserCreate500 } from '../types/PostUserCreate.ts'

/**
 * @summary Create a new user
 * {@link /user/create}
 */
export async function postUserCreate(data: PostUserCreateMutationRequest, config: Partial<RequestConfig<PostUserCreateMutationRequest>> = {}) {
  const res = await client<PostUserCreateMutationResponse, PostUserCreate500, PostUserCreateMutationRequest>({
    method: 'POST',
    url: '/user/create',
    data,
    ...config
  })

  return res.data
}