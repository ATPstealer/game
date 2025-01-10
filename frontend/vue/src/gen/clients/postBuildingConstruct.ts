import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type {
  PostBuildingConstructMutationRequest,
  PostBuildingConstructMutationResponse,
  PostBuildingConstruct401,
  PostBuildingConstruct500
} from '../types/PostBuildingConstruct.ts'

/**
 * @summary Construct a new building
 * {@link /building/construct}
 */
export async function postBuildingConstruct(
  data: PostBuildingConstructMutationRequest,
  config: Partial<RequestConfig<PostBuildingConstructMutationRequest>> = {}
) {
  const res = await client<PostBuildingConstructMutationResponse, PostBuildingConstruct401 | PostBuildingConstruct500, PostBuildingConstructMutationRequest>({
    method: 'POST',
    url: '/building/construct',
    data,
    ...config
  })
  
  return res
}