import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetEquipmentTypesQueryResponse, GetEquipmentTypes500 } from '../types/GetEquipmentTypes.ts'

/**
 * @summary Get all equipment types
 * {@link /equipment/types}
 */
export async function getEquipmentTypes(config: Partial<RequestConfig> = {}) {
  const res = await client<GetEquipmentTypesQueryResponse, GetEquipmentTypes500, unknown>({ method: 'GET', url: '/equipment/types', ...config })
  
  return res.data
}