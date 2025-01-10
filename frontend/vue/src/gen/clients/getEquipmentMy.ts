import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { GetEquipmentMyQueryResponse, GetEquipmentMyQueryParams, GetEquipmentMy401, GetEquipmentMy500 } from '../types/GetEquipmentMy.ts'

/**
 * @summary Return user's equipment
 * {@link /equipment/my}
 */
export async function getEquipmentMy(params?: GetEquipmentMyQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetEquipmentMyQueryResponse, GetEquipmentMy401 | GetEquipmentMy500, unknown>({
    method: 'GET',
    url: '/equipment/my',
    params,
    ...config
  })
  
  return res
}