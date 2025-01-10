import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type {
  PostBuildingStopWorkMutationRequest,
  PostBuildingStopWorkMutationResponse,
  PostBuildingStopWork401,
  PostBuildingStopWork500
} from '../types/PostBuildingStopWork.ts'

/**
 * @summary Stops any work in building. Later he should stop only the works available for stopping.
 * {@link /building/stop_work}
 */
export async function postBuildingStopWork(
  data: PostBuildingStopWorkMutationRequest,
  config: Partial<RequestConfig<PostBuildingStopWorkMutationRequest>> = {}
) {
  const res = await client<PostBuildingStopWorkMutationResponse, PostBuildingStopWork401 | PostBuildingStopWork500, PostBuildingStopWorkMutationRequest>({
    method: 'POST',
    url: '/building/stop_work',
    data,
    ...config
  })
  
  return res
}