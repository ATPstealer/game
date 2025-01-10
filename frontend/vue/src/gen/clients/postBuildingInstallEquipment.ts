import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type {
  PostBuildingInstallEquipmentMutationRequest,
  PostBuildingInstallEquipmentMutationResponse,
  PostBuildingInstallEquipment401,
  PostBuildingInstallEquipment500
} from '../types/PostBuildingInstallEquipment.ts'

/**
 * @summary Install equipment in a building
 * {@link /building/install_equipment}
 */
export async function postBuildingInstallEquipment(
  data: PostBuildingInstallEquipmentMutationRequest,
  config: Partial<RequestConfig<PostBuildingInstallEquipmentMutationRequest>> = {}
) {
  const res = await client<
    PostBuildingInstallEquipmentMutationResponse,
    PostBuildingInstallEquipment401 | PostBuildingInstallEquipment500,
    PostBuildingInstallEquipmentMutationRequest
  >({ method: 'POST', url: '/building/install_equipment', data, ...config })
  
  return res
}