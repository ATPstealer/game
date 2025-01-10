import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { GetMapQueryResponse, GetMap500 } from '../types/GetMap.ts'

/**
 * @description Returns the list of all map cells
 * @summary Return map cells
 * {@link /map}
 */
export async function getMap(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapQueryResponse, GetMap500, unknown>({ method: 'GET', url: '/map', ...config })
  
  return res
}