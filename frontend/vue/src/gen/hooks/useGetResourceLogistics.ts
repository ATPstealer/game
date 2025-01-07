import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetResourceLogisticsQueryResponse, GetResourceLogisticsQueryParams, GetResourceLogistics500 } from '../types/GetResourceLogistics.ts'

export const getResourceLogisticsQueryKey = (params?: MaybeRef<GetResourceLogisticsQueryParams>) =>
  [{ url: '/resource/logistics' }, ...(params ? [params] : [])] as const

export type GetResourceLogisticsQueryKey = ReturnType<typeof getResourceLogisticsQueryKey>

/**
 * @summary Get the logistics capacity in cell
 * {@link /resource/logistics}
 */
async function getResourceLogistics(params?: GetResourceLogisticsQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetResourceLogisticsQueryResponse, GetResourceLogistics500, unknown>({
    method: 'GET',
    url: '/resource/logistics',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    params,
    ...config
  })
  
  return res.data
}

export function getResourceLogisticsQueryOptions(params?: MaybeRef<GetResourceLogisticsQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getResourceLogisticsQueryKey(params)
  
  return queryOptions<GetResourceLogisticsQueryResponse, GetResourceLogistics500, GetResourceLogisticsQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getResourceLogistics(unref(params), unref(config))
    }
  })
}

/**
 * @summary Get the logistics capacity in cell
 * {@link /resource/logistics}
 */
export function useGetResourceLogistics<
  TData = GetResourceLogisticsQueryResponse,
  TQueryData = GetResourceLogisticsQueryResponse,
  TQueryKey extends QueryKey = GetResourceLogisticsQueryKey,
>(
  params?: MaybeRef<GetResourceLogisticsQueryParams>,
  options: {
    query?: Partial<QueryObserverOptions<GetResourceLogisticsQueryResponse, GetResourceLogistics500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getResourceLogisticsQueryKey(params)

  const query = useQuery({
    ...(getResourceLogisticsQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetResourceLogistics500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}