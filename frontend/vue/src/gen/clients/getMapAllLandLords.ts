import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { GetMapAllLandLordsQueryResponse, GetMapAllLandLords500 } from '../types/GetMapAllLandLords.ts'

/**
 * @summary Return all landowners
 * {@link /map/all_land_lords}
 */
export async function getMapAllLandLords(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapAllLandLordsQueryResponse, GetMapAllLandLords500, unknown>({ method: 'GET', url: '/map/all_land_lords', ...config })
  
  return res
}