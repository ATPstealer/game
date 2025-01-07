import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetMapMyQueryResponse, GetMapMy500 } from '../types/GetMapMy.ts'

export const getMapMyInfiniteQueryKey = () => [{ url: '/map/my' }] as const

export type GetMapMyInfiniteQueryKey = ReturnType<typeof getMapMyInfiniteQueryKey>

/**
 * @summary Return user's lands
 * {@link /map/my}
 */
async function getMapMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapMyQueryResponse, GetMapMy500, unknown>({
    method: 'GET',
    url: '/map/my',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getMapMyInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getMapMyInfiniteQueryKey()
  
  return infiniteQueryOptions<GetMapMyQueryResponse, GetMapMy500, GetMapMyQueryResponse, typeof queryKey, number>({
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
    query?: Partial<InfiniteQueryObserverOptions<GetMapMyQueryResponse, GetMapMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMapMyInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getMapMyInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetMapMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}