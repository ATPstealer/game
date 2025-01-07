import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type {
  DeleteBuildingDestroyMutationResponse,
  DeleteBuildingDestroyQueryParams,
  DeleteBuildingDestroy401,
  DeleteBuildingDestroy500
} from '../types/DeleteBuildingDestroy.ts'

/**
 * @summary Destroy an existing building
 * {@link /building/destroy}
 */
export async function deleteBuildingDestroy(params: DeleteBuildingDestroyQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<DeleteBuildingDestroyMutationResponse, DeleteBuildingDestroy401 | DeleteBuildingDestroy500, unknown>({
    method: 'DELETE',
    url: '/building/destroy',
    params,
    ...config
  })

  return res.data
}