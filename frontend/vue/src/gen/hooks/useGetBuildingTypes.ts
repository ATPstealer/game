import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetBuildingTypesQueryResponse, GetBuildingTypes500 } from '../types/GetBuildingTypes.ts'

export const getBuildingTypesQueryKey = () => [{ url: '/building/types' }] as const

export type GetBuildingTypesQueryKey = ReturnType<typeof getBuildingTypesQueryKey>

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

export function getBuildingTypesQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getBuildingTypesQueryKey()
  
  return queryOptions<GetBuildingTypesQueryResponse, GetBuildingTypes500, GetBuildingTypesQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getBuildingTypes(unref(config))
    }
  })
}

/**
 * @summary Get all building types
 * {@link /building/types}
 */
export function useGetBuildingTypes<
  TData = GetBuildingTypesQueryResponse,
  TQueryData = GetBuildingTypesQueryResponse,
  TQueryKey extends QueryKey = GetBuildingTypesQueryKey,
>(
  options: {
    query?: Partial<QueryObserverOptions<GetBuildingTypesQueryResponse, GetBuildingTypes500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBuildingTypesQueryKey()

  const query = useQuery({
    ...(getBuildingTypesQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetBuildingTypes500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}