import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetBuildingTypesQueryResponse, GetBuildingTypes500 } from '../types/GetBuildingTypes.ts'

/**
 * @summary Get all building types
 * {@link /building/types}
 */
export async function getBuildingTypes(config: Partial<RequestConfig> = {}) {
  const res = await client<GetBuildingTypesQueryResponse, GetBuildingTypes500, unknown>({ method: 'GET', url: '/building/types', ...config })

  return res.data
}