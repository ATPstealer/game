import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetBuildingTypesQueryResponse, GetBuildingTypes500 } from '../types/GetBuildingTypes.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getBuildingTypesQueryKey = () => [{ url: '/building/types' }] as const

export type GetBuildingTypesQueryKey = ReturnType<typeof getBuildingTypesQueryKey>

/**
 * @summary Get all building types
 * {@link /building/types}
 */
async function getBuildingTypes(config: Partial<RequestConfig> = {}) {
  const res = await client<GetBuildingTypesQueryResponse, ResponseErrorConfig<GetBuildingTypes500>, unknown>({
    method: 'GET',
    url: '/building/types',
    ...config
  })
  
  return res.data
}

export function getBuildingTypesQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getBuildingTypesQueryKey()
  
  return queryOptions<GetBuildingTypesQueryResponse, ResponseErrorConfig<GetBuildingTypes500>, GetBuildingTypesQueryResponse, typeof queryKey>({
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
    query?: Partial<QueryObserverOptions<GetBuildingTypesQueryResponse, ResponseErrorConfig<GetBuildingTypes500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBuildingTypesQueryKey()

  const query = useQuery({
    ...(getBuildingTypesQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetBuildingTypes500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}