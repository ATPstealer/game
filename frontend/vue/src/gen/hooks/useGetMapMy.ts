import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetMapMyQueryResponse, GetMapMy500 } from '../types/GetMapMy.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getMapMyQueryKey = () => [{ url: '/map/my' }] as const

export type GetMapMyQueryKey = ReturnType<typeof getMapMyQueryKey>

/**
 * @summary Return user's lands
 * {@link /map/my}
 */
async function getMapMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapMyQueryResponse, ResponseErrorConfig<GetMapMy500>, unknown>({ method: 'GET', url: '/map/my', ...config })
  
  return res.data
}

export function getMapMyQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getMapMyQueryKey()
  
  return queryOptions<GetMapMyQueryResponse, ResponseErrorConfig<GetMapMy500>, GetMapMyQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getMapMy(unref(config))
    }
  })
}

/**
 * @summary Return user's lands
 * {@link /map/my}
 */
export function useGetMapMy<TData = GetMapMyQueryResponse, TQueryData = GetMapMyQueryResponse, TQueryKey extends QueryKey = GetMapMyQueryKey>(
  options: {
    query?: Partial<QueryObserverOptions<GetMapMyQueryResponse, ResponseErrorConfig<GetMapMy500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMapMyQueryKey()

  const query = useQuery({
    ...(getMapMyQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetMapMy500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}