import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { GetBuildingMyQueryResponse, GetBuildingMyQueryParams, GetBuildingMy401, GetBuildingMy500 } from '../types/GetBuildingMy.ts'

export const getBuildingMyInfiniteQueryKey = (params?: MaybeRef<GetBuildingMyQueryParams>) => [{ url: '/building/my' }, ...(params ? [params] : [])] as const

export type GetBuildingMyInfiniteQueryKey = ReturnType<typeof getBuildingMyInfiniteQueryKey>

/**
 * @description Optionally filter by building ID.
 * @summary Fetch the user's buildings
 * {@link /building/my}
 */
async function getBuildingMy(params?: GetBuildingMyQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetBuildingMyQueryResponse, GetBuildingMy401 | GetBuildingMy500, unknown>({
    method: 'GET',
    url: '/building/my',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    params,
    ...config
  })
  
  return res.data
}

export function getBuildingMyInfiniteQueryOptions(params?: MaybeRef<GetBuildingMyQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getBuildingMyInfiniteQueryKey(params)
  
  return infiniteQueryOptions<GetBuildingMyQueryResponse, GetBuildingMy401 | GetBuildingMy500, GetBuildingMyQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      if (params) {
        params['next_page'] = pageParam as unknown as GetBuildingMyQueryParams['next_page']
      }
      
      return getBuildingMy(params, config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @description Optionally filter by building ID.
 * @summary Fetch the user's buildings
 * {@link /building/my}
 */
export function useGetBuildingMyInfinite<
  TData = InfiniteData<GetBuildingMyQueryResponse>,
  TQueryData = GetBuildingMyQueryResponse,
  TQueryKey extends QueryKey = GetBuildingMyInfiniteQueryKey,
>(
  params?: MaybeRef<GetBuildingMyQueryParams>,
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetBuildingMyQueryResponse, GetBuildingMy401 | GetBuildingMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBuildingMyInfiniteQueryKey(params)

  const query = useInfiniteQuery({
    ...(getBuildingMyInfiniteQueryOptions(params, config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetBuildingMy401 | GetBuildingMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}