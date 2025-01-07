import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetStorageMyQueryResponse, GetStorageMy401, GetStorageMy500 } from '../types/GetStorageMy.ts'

export const getStorageMyInfiniteQueryKey = () => [{ url: '/storage/my' }] as const

export type GetStorageMyInfiniteQueryKey = ReturnType<typeof getStorageMyInfiniteQueryKey>

/**
 * @summary Return user's storages
 * {@link /storage/my}
 */
async function getStorageMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetStorageMyQueryResponse, GetStorageMy401 | GetStorageMy500, unknown>({
    method: 'GET',
    url: '/storage/my',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getStorageMyInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getStorageMyInfiniteQueryKey()
  
  return infiniteQueryOptions<GetStorageMyQueryResponse, GetStorageMy401 | GetStorageMy500, GetStorageMyQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getStorageMy(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Return user's storages
 * {@link /storage/my}
 */
export function useGetStorageMyInfinite<
  TData = InfiniteData<GetStorageMyQueryResponse>,
  TQueryData = GetStorageMyQueryResponse,
  TQueryKey extends QueryKey = GetStorageMyInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetStorageMyQueryResponse, GetStorageMy401 | GetStorageMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getStorageMyInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getStorageMyInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetStorageMy401 | GetStorageMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}