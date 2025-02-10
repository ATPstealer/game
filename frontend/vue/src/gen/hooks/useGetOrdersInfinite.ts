import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { GetOrdersQueryResponse, GetOrdersQueryParams, GetOrders500 } from '../types/GetOrders.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const getOrdersInfiniteQueryKey = (params?: MaybeRef<GetOrdersQueryParams>) => [{ url: '/orders' }, ...(params ? [params] : [])] as const

export type GetOrdersInfiniteQueryKey = ReturnType<typeof getOrdersInfiniteQueryKey>

/**
 * @summary Fetches orders based on various query parameters
 * {@link /orders}
 */
async function getOrders(params?: GetOrdersQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetOrdersQueryResponse, ResponseErrorConfig<GetOrders500>, unknown>({ method: 'GET', url: '/orders', params, ...config })
  
  return res.data
}

export function getOrdersInfiniteQueryOptions(params?: MaybeRef<GetOrdersQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getOrdersInfiniteQueryKey(params)
  
  return infiniteQueryOptions<GetOrdersQueryResponse, ResponseErrorConfig<GetOrders500>, GetOrdersQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      if (params) {
        params['next_page'] = pageParam as unknown as GetOrdersQueryParams['next_page']
      }
      
      return getOrders(params, config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Fetches orders based on various query parameters
 * {@link /orders}
 */
export function useGetOrdersInfinite<
  TData = InfiniteData<GetOrdersQueryResponse>,
  TQueryData = GetOrdersQueryResponse,
  TQueryKey extends QueryKey = GetOrdersInfiniteQueryKey,
>(
  params?: MaybeRef<GetOrdersQueryParams>,
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetOrdersQueryResponse, ResponseErrorConfig<GetOrders500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getOrdersInfiniteQueryKey(params)

  const query = useInfiniteQuery({
    ...(getOrdersInfiniteQueryOptions(params, config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, ResponseErrorConfig<GetOrders500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}