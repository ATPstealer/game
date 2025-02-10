import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetSettingsQueryResponse, GetSettings500 } from '../types/GetSettings.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getSettingsQueryKey = () => [{ url: '/settings' }] as const

export type GetSettingsQueryKey = ReturnType<typeof getSettingsQueryKey>

/**
 * @description X Y dimension, Interest rate, etc
 * @summary Get General Game Settings
 * {@link /settings}
 */
async function getSettings(config: Partial<RequestConfig> = {}) {
  const res = await client<GetSettingsQueryResponse, ResponseErrorConfig<GetSettings500>, unknown>({ method: 'GET', url: '/settings', ...config })
  
  return res.data
}

export function getSettingsQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getSettingsQueryKey()
  
  return queryOptions<GetSettingsQueryResponse, ResponseErrorConfig<GetSettings500>, GetSettingsQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getSettings(unref(config))
    }
  })
}

/**
 * @description X Y dimension, Interest rate, etc
 * @summary Get General Game Settings
 * {@link /settings}
 */
export function useGetSettings<TData = GetSettingsQueryResponse, TQueryData = GetSettingsQueryResponse, TQueryKey extends QueryKey = GetSettingsQueryKey>(
  options: {
    query?: Partial<QueryObserverOptions<GetSettingsQueryResponse, ResponseErrorConfig<GetSettings500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getSettingsQueryKey()

  const query = useQuery({
    ...(getSettingsQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetSettings500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}