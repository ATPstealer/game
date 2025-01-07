import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetBuildingMyQueryResponse, GetBuildingMyQueryParams, GetBuildingMy401, GetBuildingMy500 } from '../types/GetBuildingMy.ts'

export const getBuildingMyQueryKey = (params?: MaybeRef<GetBuildingMyQueryParams>) => [{ url: '/building/my' }, ...(params ? [params] : [])] as const

export type GetBuildingMyQueryKey = ReturnType<typeof getBuildingMyQueryKey>

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

export function getBuildingMyQueryOptions(params?: MaybeRef<GetBuildingMyQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getBuildingMyQueryKey(params)
  
  return queryOptions<GetBuildingMyQueryResponse, GetBuildingMy401 | GetBuildingMy500, GetBuildingMyQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getBuildingMy(unref(params), unref(config))
    }
  })
}

/**
 * @description Optionally filter by building ID.
 * @summary Fetch the user's buildings
 * {@link /building/my}
 */
export function useGetBuildingMy<
  TData = GetBuildingMyQueryResponse,
  TQueryData = GetBuildingMyQueryResponse,
  TQueryKey extends QueryKey = GetBuildingMyQueryKey,
>(
  params?: MaybeRef<GetBuildingMyQueryParams>,
  options: {
    query?: Partial<QueryObserverOptions<GetBuildingMyQueryResponse, GetBuildingMy401 | GetBuildingMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBuildingMyQueryKey(params)

  const query = useQuery({
    ...(getBuildingMyQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetBuildingMy401 | GetBuildingMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}