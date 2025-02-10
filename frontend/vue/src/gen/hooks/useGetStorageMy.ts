import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetStorageMyQueryResponse, GetStorageMy401, GetStorageMy500 } from '../types/GetStorageMy.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getStorageMyQueryKey = () => [{ url: '/storage/my' }] as const

export type GetStorageMyQueryKey = ReturnType<typeof getStorageMyQueryKey>

/**
 * @summary Return user's storages
 * {@link /storage/my}
 */
async function getStorageMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetStorageMyQueryResponse, ResponseErrorConfig<GetStorageMy401 | GetStorageMy500>, unknown>({
    method: 'GET',
    url: '/storage/my',
    ...config
  })
  
  return res.data
}

export function getStorageMyQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getStorageMyQueryKey()
  
  return queryOptions<GetStorageMyQueryResponse, ResponseErrorConfig<GetStorageMy401 | GetStorageMy500>, GetStorageMyQueryResponse, typeof queryKey>({
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
    query?: Partial<QueryObserverOptions<GetStorageMyQueryResponse, ResponseErrorConfig<GetStorageMy401 | GetStorageMy500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getStorageMyQueryKey()

  const query = useQuery({
    ...(getStorageMyQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetStorageMy401 | GetStorageMy500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}