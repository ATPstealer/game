import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type {
  PostBuildingEmergencyHiringMutationRequest,
  PostBuildingEmergencyHiringMutationResponse,
  PostBuildingEmergencyHiring401,
  PostBuildingEmergencyHiring500
} from '../types/PostBuildingEmergencyHiring.ts'

/**
 * @summary Expensive fast hiring
 * {@link /building/emergency_hiring}
 */
export async function postBuildingEmergencyHiring(
  data: PostBuildingEmergencyHiringMutationRequest,
  config: Partial<RequestConfig<PostBuildingEmergencyHiringMutationRequest>> = {}
) {
  const res = await client<
    PostBuildingEmergencyHiringMutationResponse,
    PostBuildingEmergencyHiring401 | PostBuildingEmergencyHiring500,
    PostBuildingEmergencyHiringMutationRequest
  >({ method: 'POST', url: '/building/emergency_hiring', data, ...config })
  
  return res.data
}