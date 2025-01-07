import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { GetResourceLogisticsQueryResponse, GetResourceLogisticsQueryParams, GetResourceLogistics500 } from '../types/GetResourceLogistics.ts'

export const getResourceLogisticsInfiniteQueryKey = (params?: MaybeRef<GetResourceLogisticsQueryParams>) =>
  [{ url: '/resource/logistics' }, ...(params ? [params] : [])] as const

export type GetResourceLogisticsInfiniteQueryKey = ReturnType<typeof getResourceLogisticsInfiniteQueryKey>

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

export function getResourceLogisticsInfiniteQueryOptions(params?: MaybeRef<GetResourceLogisticsQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getResourceLogisticsInfiniteQueryKey(params)
  
  return infiniteQueryOptions<GetResourceLogisticsQueryResponse, GetResourceLogistics500, GetResourceLogisticsQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      if (params) {
        params['next_page'] = pageParam as unknown as GetResourceLogisticsQueryParams['next_page']
      }
      
      return getResourceLogistics(params, config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Get the logistics capacity in cell
 * {@link /resource/logistics}
 */
export function useGetResourceLogisticsInfinite<
  TData = InfiniteData<GetResourceLogisticsQueryResponse>,
  TQueryData = GetResourceLogisticsQueryResponse,
  TQueryKey extends QueryKey = GetResourceLogisticsInfiniteQueryKey,
>(
  params?: MaybeRef<GetResourceLogisticsQueryParams>,
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetResourceLogisticsQueryResponse, GetResourceLogistics500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getResourceLogisticsInfiniteQueryKey(params)

  const query = useInfiniteQuery({
    ...(getResourceLogisticsInfiniteQueryOptions(params, config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetResourceLogistics500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}