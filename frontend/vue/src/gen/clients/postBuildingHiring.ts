import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type {
  PostBuildingHiringMutationRequest,
  PostBuildingHiringMutationResponse,
  PostBuildingHiring401,
  PostBuildingHiring500
} from '../types/PostBuildingHiring.ts'

/**
 * @summary Set hiring details for a building
 * {@link /building/hiring}
 */
export async function postBuildingHiring(data?: PostBuildingHiringMutationRequest, config: Partial<RequestConfig<PostBuildingHiringMutationRequest>> = {}) {
  const res = await client<PostBuildingHiringMutationResponse, PostBuildingHiring401 | PostBuildingHiring500, PostBuildingHiringMutationRequest>({
    method: 'POST',
    url: '/building/hiring',
    data,
    ...config
  })
  
  return res.data
}