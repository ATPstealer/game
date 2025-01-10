import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetMapQueryResponse, GetMap500 } from '../types/GetMap.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const getMapInfiniteQueryKey = () => [{ url: '/map' }] as const

export type GetMapInfiniteQueryKey = ReturnType<typeof getMapInfiniteQueryKey>

/**
 * @description Returns the list of all map cells
 * @summary Return map cells
 * {@link /map}
 */
async function getMap(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapQueryResponse, GetMap500, unknown>({
    method: 'GET',
    url: '/map',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getMapInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getMapInfiniteQueryKey()
  
  return infiniteQueryOptions<GetMapQueryResponse, GetMap500, GetMapQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getMap(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @description Returns the list of all map cells
 * @summary Return map cells
 * {@link /map}
 */
export function useGetMapInfinite<
  TData = InfiniteData<GetMapQueryResponse>,
  TQueryData = GetMapQueryResponse,
  TQueryKey extends QueryKey = GetMapInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetMapQueryResponse, GetMap500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMapInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getMapInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetMap500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}