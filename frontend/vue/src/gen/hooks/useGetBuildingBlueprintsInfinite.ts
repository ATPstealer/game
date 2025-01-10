import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { GetBuildingBlueprintsQueryResponse, GetBuildingBlueprintsQueryParams, GetBuildingBlueprints500 } from '../types/GetBuildingBlueprints.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const getBuildingBlueprintsInfiniteQueryKey = (params?: MaybeRef<GetBuildingBlueprintsQueryParams>) =>
  [{ url: '/building/blueprints' }, ...(params ? [params] : [])] as const

export type GetBuildingBlueprintsInfiniteQueryKey = ReturnType<typeof getBuildingBlueprintsInfiniteQueryKey>

/**
 * @description Fetches a list of blueprints. If an 'id' query parameter is provided, fetches the blueprint with the specified ID.
 * @summary Get blueprints
 * {@link /building/blueprints}
 */
async function getBuildingBlueprints(params?: GetBuildingBlueprintsQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetBuildingBlueprintsQueryResponse, GetBuildingBlueprints500, unknown>({
    method: 'GET',
    url: '/building/blueprints',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    params,
    ...config
  })
  
  return res.data
}

export function getBuildingBlueprintsInfiniteQueryOptions(params?: MaybeRef<GetBuildingBlueprintsQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getBuildingBlueprintsInfiniteQueryKey(params)
  
  return infiniteQueryOptions<GetBuildingBlueprintsQueryResponse, GetBuildingBlueprints500, GetBuildingBlueprintsQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      if (params) {
        params['next_page'] = pageParam as unknown as GetBuildingBlueprintsQueryParams['next_page']
      }
      
      return getBuildingBlueprints(params, config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @description Fetches a list of blueprints. If an 'id' query parameter is provided, fetches the blueprint with the specified ID.
 * @summary Get blueprints
 * {@link /building/blueprints}
 */
export function useGetBuildingBlueprintsInfinite<
  TData = InfiniteData<GetBuildingBlueprintsQueryResponse>,
  TQueryData = GetBuildingBlueprintsQueryResponse,
  TQueryKey extends QueryKey = GetBuildingBlueprintsInfiniteQueryKey,
>(
  params?: MaybeRef<GetBuildingBlueprintsQueryParams>,
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetBuildingBlueprintsQueryResponse, GetBuildingBlueprints500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBuildingBlueprintsInfiniteQueryKey(params)

  const query = useInfiniteQuery({
    ...(getBuildingBlueprintsInfiniteQueryOptions(params, config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetBuildingBlueprints500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}