import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetMapAllLandLordsQueryResponse, GetMapAllLandLords500 } from '../types/GetMapAllLandLords.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const getMapAllLandLordsInfiniteQueryKey = () => [{ url: '/map/all_land_lords' }] as const

export type GetMapAllLandLordsInfiniteQueryKey = ReturnType<typeof getMapAllLandLordsInfiniteQueryKey>

/**
 * @summary Return all landowners
 * {@link /map/all_land_lords}
 */
async function getMapAllLandLords(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapAllLandLordsQueryResponse, ResponseErrorConfig<GetMapAllLandLords500>, unknown>({
    method: 'GET',
    url: '/map/all_land_lords',
    ...config
  })
  
  return res.data
}

export function getMapAllLandLordsInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getMapAllLandLordsInfiniteQueryKey()
  
  return infiniteQueryOptions<
    GetMapAllLandLordsQueryResponse,
    ResponseErrorConfig<GetMapAllLandLords500>,
    GetMapAllLandLordsQueryResponse,
    typeof queryKey,
    number
  >({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getMapAllLandLords(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Return all landowners
 * {@link /map/all_land_lords}
 */
export function useGetMapAllLandLordsInfinite<
  TData = InfiniteData<GetMapAllLandLordsQueryResponse>,
  TQueryData = GetMapAllLandLordsQueryResponse,
  TQueryKey extends QueryKey = GetMapAllLandLordsInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetMapAllLandLordsQueryResponse, ResponseErrorConfig<GetMapAllLandLords500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMapAllLandLordsInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getMapAllLandLordsInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, ResponseErrorConfig<GetMapAllLandLords500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}