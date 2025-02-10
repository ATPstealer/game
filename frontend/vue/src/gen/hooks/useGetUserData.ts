import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetUserDataQueryResponse, GetUserData401 } from '../types/GetUserData.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getUserDataQueryKey = () => [{ url: '/user/data' }] as const

export type GetUserDataQueryKey = ReturnType<typeof getUserDataQueryKey>

/**
 * @summary Get user data
 * {@link /user/data}
 */
async function getUserData(config: Partial<RequestConfig> = {}) {
  const res = await client<GetUserDataQueryResponse, ResponseErrorConfig<GetUserData401>, unknown>({ method: 'GET', url: '/user/data', ...config })
  
  return res.data
}

export function getUserDataQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getUserDataQueryKey()
  
  return queryOptions<GetUserDataQueryResponse, ResponseErrorConfig<GetUserData401>, GetUserDataQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getUserData(unref(config))
    }
  })
}

/**
 * @summary Get user data
 * {@link /user/data}
 */
export function useGetUserData<TData = GetUserDataQueryResponse, TQueryData = GetUserDataQueryResponse, TQueryKey extends QueryKey = GetUserDataQueryKey>(
  options: {
    query?: Partial<QueryObserverOptions<GetUserDataQueryResponse, ResponseErrorConfig<GetUserData401>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getUserDataQueryKey()

  const query = useQuery({
    ...(getUserDataQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetUserData401>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}