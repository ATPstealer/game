import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { GetBuildingBlueprintsQueryResponse, GetBuildingBlueprintsQueryParams, GetBuildingBlueprints500 } from '../types/GetBuildingBlueprints.ts'

/**
 * @description Fetches a list of blueprints. If an 'id' query parameter is provided, fetches the blueprint with the specified ID.
 * @summary Get blueprints
 * {@link /building/blueprints}
 */
export async function getBuildingBlueprints(params?: GetBuildingBlueprintsQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetBuildingBlueprintsQueryResponse, GetBuildingBlueprints500, unknown>({
    method: 'GET',
    url: '/building/blueprints',
    params,
    ...config
  })
  
  return res
}