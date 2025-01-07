import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetMapMyQueryResponse, GetMapMy500 } from '../types/GetMapMy.ts'

/**
 * @summary Return user's lands
 * {@link /map/my}
 */
export async function getMapMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapMyQueryResponse, GetMapMy500, unknown>({ method: 'GET', url: '/map/my', ...config })

  return res.data
}