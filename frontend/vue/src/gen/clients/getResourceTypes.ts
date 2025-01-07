import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetResourceTypesQueryResponse, GetResourceTypes500 } from '../types/GetResourceTypes.ts'

/**
 * @summary Return all resource types from database
 * {@link /resource/types}
 */
export async function getResourceTypes(config: Partial<RequestConfig> = {}) {
  const res = await client<GetResourceTypesQueryResponse, GetResourceTypes500, unknown>({ method: 'GET', url: '/resource/types', ...config })
  
  return res.data
}