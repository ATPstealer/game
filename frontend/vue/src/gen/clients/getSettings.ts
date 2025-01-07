import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type { GetSettingsQueryResponse, GetSettings500 } from '../types/GetSettings.ts'

/**
 * @description X Y dimension, Interest rate, etc
 * @summary Get General Game Settings
 * {@link /settings}
 */
export async function getSettings(config: Partial<RequestConfig> = {}) {
  const res = await client<GetSettingsQueryResponse, GetSettings500, unknown>({ method: 'GET', url: '/settings', ...config })

  return res.data
}