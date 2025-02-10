import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetMapMyQueryResponse, GetMapMy500 } from '../types/GetMapMy.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const getMapMyInfiniteQueryKey = () => [{ url: '/map/my' }] as const

export type GetMapMyInfiniteQueryKey = ReturnType<typeof getMapMyInfiniteQueryKey>

/**
 * @summary Return user's lands
 * {@link /map/my}
 */
async function getMapMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapMyQueryResponse, ResponseErrorConfig<GetMapMy500>, unknown>({ method: 'GET', url: '/map/my', ...config })
  
  return res.data
}

export function getMapMyInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getMapMyInfiniteQueryKey()
  
  return infiniteQueryOptions<GetMapMyQueryResponse, ResponseErrorConfig<GetMapMy500>, GetMapMyQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getMapMy(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Return user's lands
 * {@link /map/my}
 */
export function useGetMapMyInfinite<
  TData = InfiniteData<GetMapMyQueryResponse>,
  TQueryData = GetMapMyQueryResponse,
  TQueryKey extends QueryKey = GetMapMyInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetMapMyQueryResponse, ResponseErrorConfig<GetMapMy500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMapMyInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getMapMyInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, ResponseErrorConfig<GetMapMy500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}