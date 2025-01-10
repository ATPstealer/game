import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type { GetStorageMyQueryResponse, GetStorageMy401, GetStorageMy500 } from '../types/GetStorageMy.ts'

/**
 * @summary Return user's storages
 * {@link /storage/my}
 */
export async function getStorageMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetStorageMyQueryResponse, GetStorageMy401 | GetStorageMy500, unknown>({ method: 'GET', url: '/storage/my', ...config })
  
  return res
}