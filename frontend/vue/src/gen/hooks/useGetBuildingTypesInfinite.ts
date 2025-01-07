import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetBuildingTypesQueryResponse, GetBuildingTypes500 } from '../types/GetBuildingTypes.ts'

export const getBuildingTypesInfiniteQueryKey = () => [{ url: '/building/types' }] as const

export type GetBuildingTypesInfiniteQueryKey = ReturnType<typeof getBuildingTypesInfiniteQueryKey>

/**
 * @summary Get all building types
 * {@link /building/types}
 */
async function getBuildingTypes(config: Partial<RequestConfig> = {}) {
  const res = await client<GetBuildingTypesQueryResponse, GetBuildingTypes500, unknown>({
    method: 'GET',
    url: '/building/types',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getBuildingTypesInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getBuildingTypesInfiniteQueryKey()
  
  return infiniteQueryOptions<GetBuildingTypesQueryResponse, GetBuildingTypes500, GetBuildingTypesQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getBuildingTypes(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Get all building types
 * {@link /building/types}
 */
export function useGetBuildingTypesInfinite<
  TData = InfiniteData<GetBuildingTypesQueryResponse>,
  TQueryData = GetBuildingTypesQueryResponse,
  TQueryKey extends QueryKey = GetBuildingTypesInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetBuildingTypesQueryResponse, GetBuildingTypes500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBuildingTypesInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getBuildingTypesInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetBuildingTypes500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}