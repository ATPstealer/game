import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetBuildingBlueprintsQueryResponse, GetBuildingBlueprintsQueryParams, GetBuildingBlueprints500 } from '../types/GetBuildingBlueprints.ts'

export const getBuildingBlueprintsQueryKey = (params?: MaybeRef<GetBuildingBlueprintsQueryParams>) =>
  [{ url: '/building/blueprints' }, ...(params ? [params] : [])] as const

export type GetBuildingBlueprintsQueryKey = ReturnType<typeof getBuildingBlueprintsQueryKey>

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

export function getBuildingBlueprintsQueryOptions(params?: MaybeRef<GetBuildingBlueprintsQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getBuildingBlueprintsQueryKey(params)
  
  return queryOptions<GetBuildingBlueprintsQueryResponse, GetBuildingBlueprints500, GetBuildingBlueprintsQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getBuildingBlueprints(unref(params), unref(config))
    }
  })
}

/**
 * @description Fetches a list of blueprints. If an 'id' query parameter is provided, fetches the blueprint with the specified ID.
 * @summary Get blueprints
 * {@link /building/blueprints}
 */
export function useGetBuildingBlueprints<
  TData = GetBuildingBlueprintsQueryResponse,
  TQueryData = GetBuildingBlueprintsQueryResponse,
  TQueryKey extends QueryKey = GetBuildingBlueprintsQueryKey,
>(
  params?: MaybeRef<GetBuildingBlueprintsQueryParams>,
  options: {
    query?: Partial<QueryObserverOptions<GetBuildingBlueprintsQueryResponse, GetBuildingBlueprints500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBuildingBlueprintsQueryKey(params)

  const query = useQuery({
    ...(getBuildingBlueprintsQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetBuildingBlueprints500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}