import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetResourceMyLogisticsQueryResponse, GetResourceMyLogistics401, GetResourceMyLogistics500 } from '../types/GetResourceMyLogistics.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const getResourceMyLogisticsInfiniteQueryKey = () => [{ url: '/resource/my_logistics' }] as const

export type GetResourceMyLogisticsInfiniteQueryKey = ReturnType<typeof getResourceMyLogisticsInfiniteQueryKey>

/**
 * @summary Get user's logistics tasks
 * {@link /resource/my_logistics}
 */
async function getResourceMyLogistics(config: Partial<RequestConfig> = {}) {
  const res = await client<GetResourceMyLogisticsQueryResponse, ResponseErrorConfig<GetResourceMyLogistics401 | GetResourceMyLogistics500>, unknown>({
    method: 'GET',
    url: '/resource/my_logistics',
    ...config
  })
  
  return res.data
}

export function getResourceMyLogisticsInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getResourceMyLogisticsInfiniteQueryKey()
  
  return infiniteQueryOptions<
    GetResourceMyLogisticsQueryResponse,
    ResponseErrorConfig<GetResourceMyLogistics401 | GetResourceMyLogistics500>,
    GetResourceMyLogisticsQueryResponse,
    typeof queryKey,
    number
  >({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getResourceMyLogistics(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Get user's logistics tasks
 * {@link /resource/my_logistics}
 */
export function useGetResourceMyLogisticsInfinite<
  TData = InfiniteData<GetResourceMyLogisticsQueryResponse>,
  TQueryData = GetResourceMyLogisticsQueryResponse,
  TQueryKey extends QueryKey = GetResourceMyLogisticsInfiniteQueryKey,
>(
  options: {
    query?: Partial<
      InfiniteQueryObserverOptions<
        GetResourceMyLogisticsQueryResponse,
        ResponseErrorConfig<GetResourceMyLogistics401 | GetResourceMyLogistics500>,
        TData,
        TQueryData,
        TQueryKey
      >
    >;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getResourceMyLogisticsInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getResourceMyLogisticsInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, ResponseErrorConfig<GetResourceMyLogistics401 | GetResourceMyLogistics500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}