import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type {
  PostBuildingStartWorkMutationRequest,
  PostBuildingStartWorkMutationResponse,
  PostBuildingStartWork401,
  PostBuildingStartWork500
} from '../types/PostBuildingStartWork.ts'

/**
 * @summary Start work in the building
 * {@link /building/start_work}
 */
export async function postBuildingStartWork(
  data: PostBuildingStartWorkMutationRequest,
  config: Partial<RequestConfig<PostBuildingStartWorkMutationRequest>> = {}
) {
  const res = await client<PostBuildingStartWorkMutationResponse, PostBuildingStartWork401 | PostBuildingStartWork500, PostBuildingStartWorkMutationRequest>({
    method: 'POST',
    url: '/building/start_work',
    data,
    ...config
  })
  
  return res
}