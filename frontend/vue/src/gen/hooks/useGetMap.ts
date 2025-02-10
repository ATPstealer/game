import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetMapQueryResponse, GetMap500 } from '../types/GetMap.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getMapQueryKey = () => [{ url: '/map/' }] as const

export type GetMapQueryKey = ReturnType<typeof getMapQueryKey>

/**
 * @description Returns the list of all map cells
 * @summary Return map cells
 * {@link /map}
 */
async function getMap(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapQueryResponse, ResponseErrorConfig<GetMap500>, unknown>({ method: 'GET', url: '/map/', ...config })
  
  return res.data
}

export function getMapQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getMapQueryKey()
  
  return queryOptions<GetMapQueryResponse, ResponseErrorConfig<GetMap500>, GetMapQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getMap(unref(config))
    }
  })
}

/**
 * @description Returns the list of all map cells
 * @summary Return map cells
 * {@link /map}
 */
export function useGetMap<TData = GetMapQueryResponse, TQueryData = GetMapQueryResponse, TQueryKey extends QueryKey = GetMapQueryKey>(
  options: {
    query?: Partial<QueryObserverOptions<GetMapQueryResponse, ResponseErrorConfig<GetMap500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMapQueryKey()

  const query = useQuery({
    ...(getMapQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetMap500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}