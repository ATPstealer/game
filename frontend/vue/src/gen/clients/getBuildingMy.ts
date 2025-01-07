import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetBuildingMyQueryResponse, GetBuildingMyQueryParams, GetBuildingMy401, GetBuildingMy500 } from '../types/GetBuildingMy.ts'

/**
 * @description Optionally filter by building ID.
 * @summary Fetch the user's buildings
 * {@link /building/my}
 */
export async function getBuildingMy(params?: GetBuildingMyQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetBuildingMyQueryResponse, GetBuildingMy401 | GetBuildingMy500, unknown>({ method: 'GET', url: '/building/my', params, ...config })
  
  return res.data
}