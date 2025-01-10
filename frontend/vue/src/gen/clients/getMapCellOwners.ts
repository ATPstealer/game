import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { GetMapCellOwnersQueryResponse, GetMapCellOwnersQueryParams, GetMapCellOwners500 } from '../types/GetMapCellOwners.ts'

/**
 * @summary Get the landlords in cell
 * {@link /map/cell_owners}
 */
export async function getMapCellOwners(params: GetMapCellOwnersQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapCellOwnersQueryResponse, GetMapCellOwners500, unknown>({ method: 'GET', url: '/map/cell_owners', params, ...config })
  
  return res
}