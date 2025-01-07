import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetMapAllLandLordsQueryResponse, GetMapAllLandLords500 } from '../types/GetMapAllLandLords.ts'

export const getMapAllLandLordsQueryKey = () => [{ url: '/map/all_land_lords' }] as const

export type GetMapAllLandLordsQueryKey = ReturnType<typeof getMapAllLandLordsQueryKey>

/**
 * @summary Return all landowners
 * {@link /map/all_land_lords}
 */
async function getMapAllLandLords(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapAllLandLordsQueryResponse, GetMapAllLandLords500, unknown>({
    method: 'GET',
    url: '/map/all_land_lords',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getMapAllLandLordsQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getMapAllLandLordsQueryKey()
  
  return queryOptions<GetMapAllLandLordsQueryResponse, GetMapAllLandLords500, GetMapAllLandLordsQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getMapAllLandLords(unref(config))
    }
  })
}

/**
 * @summary Return all landowners
 * {@link /map/all_land_lords}
 */
export function useGetMapAllLandLords<
  TData = GetMapAllLandLordsQueryResponse,
  TQueryData = GetMapAllLandLordsQueryResponse,
  TQueryKey extends QueryKey = GetMapAllLandLordsQueryKey,
>(
  options: {
    query?: Partial<QueryObserverOptions<GetMapAllLandLordsQueryResponse, GetMapAllLandLords500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMapAllLandLordsQueryKey()

  const query = useQuery({
    ...(getMapAllLandLordsQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetMapAllLandLords500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}