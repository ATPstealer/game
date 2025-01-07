import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetResourceMyLogisticsQueryResponse, GetResourceMyLogistics401, GetResourceMyLogistics500 } from '../types/GetResourceMyLogistics.ts'

export const getResourceMyLogisticsQueryKey = () => [{ url: '/resource/my_logistics' }] as const

export type GetResourceMyLogisticsQueryKey = ReturnType<typeof getResourceMyLogisticsQueryKey>

/**
 * @summary Get user's logistics tasks
 * {@link /resource/my_logistics}
 */
async function getResourceMyLogistics(config: Partial<RequestConfig> = {}) {
  const res = await client<GetResourceMyLogisticsQueryResponse, GetResourceMyLogistics401 | GetResourceMyLogistics500, unknown>({
    method: 'GET',
    url: '/resource/my_logistics',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getResourceMyLogisticsQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getResourceMyLogisticsQueryKey()
  
  return queryOptions<
    GetResourceMyLogisticsQueryResponse,
    GetResourceMyLogistics401 | GetResourceMyLogistics500,
    GetResourceMyLogisticsQueryResponse,
    typeof queryKey
  >({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getResourceMyLogistics(unref(config))
    }
  })
}

/**
 * @summary Get user's logistics tasks
 * {@link /resource/my_logistics}
 */
export function useGetResourceMyLogistics<
  TData = GetResourceMyLogisticsQueryResponse,
  TQueryData = GetResourceMyLogisticsQueryResponse,
  TQueryKey extends QueryKey = GetResourceMyLogisticsQueryKey,
>(
  options: {
    query?: Partial<
      QueryObserverOptions<GetResourceMyLogisticsQueryResponse, GetResourceMyLogistics401 | GetResourceMyLogistics500, TData, TQueryData, TQueryKey>
    >;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getResourceMyLogisticsQueryKey()

  const query = useQuery({
    ...(getResourceMyLogisticsQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetResourceMyLogistics401 | GetResourceMyLogistics500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}