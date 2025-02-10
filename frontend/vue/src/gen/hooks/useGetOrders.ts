import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetOrdersQueryResponse, GetOrdersQueryParams, GetOrders500 } from '../types/GetOrders.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getOrdersQueryKey = (params?: MaybeRef<GetOrdersQueryParams>) => [{ url: '/orders' }, ...(params ? [params] : [])] as const

export type GetOrdersQueryKey = ReturnType<typeof getOrdersQueryKey>

/**
 * @summary Fetches orders based on various query parameters
 * {@link /orders}
 */
async function getOrders(params?: GetOrdersQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetOrdersQueryResponse, ResponseErrorConfig<GetOrders500>, unknown>({ method: 'GET', url: '/orders', params, ...config })
  
  return res.data
}

export function getOrdersQueryOptions(params?: MaybeRef<GetOrdersQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getOrdersQueryKey(params)
  
  return queryOptions<GetOrdersQueryResponse, ResponseErrorConfig<GetOrders500>, GetOrdersQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getOrders(unref(params), unref(config))
    }
  })
}

/**
 * @summary Fetches orders based on various query parameters
 * {@link /orders}
 */
export function useGetOrders<TData = GetOrdersQueryResponse, TQueryData = GetOrdersQueryResponse, TQueryKey extends QueryKey = GetOrdersQueryKey>(
  params?: MaybeRef<GetOrdersQueryParams>,
  options: {
    query?: Partial<QueryObserverOptions<GetOrdersQueryResponse, ResponseErrorConfig<GetOrders500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getOrdersQueryKey(params)

  const query = useQuery({
    ...(getOrdersQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetOrders500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}