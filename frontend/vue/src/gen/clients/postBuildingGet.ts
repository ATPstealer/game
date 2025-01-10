import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { PostBuildingGetMutationRequest, PostBuildingGetMutationResponse, PostBuildingGet500 } from '../types/PostBuildingGet.ts'

/**
 * @summary Fetch the list of buildings
 * {@link /building/get}
 */
export async function postBuildingGet(data?: PostBuildingGetMutationRequest, config: Partial<RequestConfig<PostBuildingGetMutationRequest>> = {}) {
  const res = await client<PostBuildingGetMutationResponse, PostBuildingGet500, PostBuildingGetMutationRequest>({
    method: 'POST',
    url: '/building/get',
    data,
    ...config
  })
  
  return res
}