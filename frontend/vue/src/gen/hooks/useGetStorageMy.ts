import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetStorageMyQueryResponse, GetStorageMy401, GetStorageMy500 } from '../types/GetStorageMy.ts'

export const getStorageMyQueryKey = () => [{ url: '/storage/my' }] as const

export type GetStorageMyQueryKey = ReturnType<typeof getStorageMyQueryKey>

/**
 * @summary Return user's storages
 * {@link /storage/my}
 */
async function getStorageMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetStorageMyQueryResponse, GetStorageMy401 | GetStorageMy500, unknown>({
    method: 'GET',
    url: '/storage/my',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getStorageMyQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getStorageMyQueryKey()
  
  return queryOptions<GetStorageMyQueryResponse, GetStorageMy401 | GetStorageMy500, GetStorageMyQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getStorageMy(unref(config))
    }
  })
}

/**
 * @summary Return user's storages
 * {@link /storage/my}
 */
export function useGetStorageMy<TData = GetStorageMyQueryResponse, TQueryData = GetStorageMyQueryResponse, TQueryKey extends QueryKey = GetStorageMyQueryKey>(
  options: {
    query?: Partial<QueryObserverOptions<GetStorageMyQueryResponse, GetStorageMy401 | GetStorageMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getStorageMyQueryKey()

  const query = useQuery({
    ...(getStorageMyQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetStorageMy401 | GetStorageMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}