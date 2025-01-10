import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetResourceTypesQueryResponse, GetResourceTypes500 } from '../types/GetResourceTypes.ts'
import type { RequestConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getResourceTypesQueryKey = () => [{ url: '/resource/types' }] as const

export type GetResourceTypesQueryKey = ReturnType<typeof getResourceTypesQueryKey>

/**
 * @summary Return all resource types from database
 * {@link /resource/types}
 */
async function getResourceTypes(config: Partial<RequestConfig> = {}) {
  const res = await client<GetResourceTypesQueryResponse, GetResourceTypes500, unknown>({
    method: 'GET',
    url: '/resource/types',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getResourceTypesQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getResourceTypesQueryKey()
  
  return queryOptions<GetResourceTypesQueryResponse, GetResourceTypes500, GetResourceTypesQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getResourceTypes(unref(config))
    }
  })
}

/**
 * @summary Return all resource types from database
 * {@link /resource/types}
 */
export function useGetResourceTypes<
  TData = GetResourceTypesQueryResponse,
  TQueryData = GetResourceTypesQueryResponse,
  TQueryKey extends QueryKey = GetResourceTypesQueryKey,
>(
  options: {
    query?: Partial<QueryObserverOptions<GetResourceTypesQueryResponse, GetResourceTypes500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getResourceTypesQueryKey()

  const query = useQuery({
    ...(getResourceTypesQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetResourceTypes500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}