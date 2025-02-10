import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetUsernamesByPrefixQueryResponse, GetUsernamesByPrefixQueryParams } from '../types/GetUsernamesByPrefix.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getUsernamesByPrefixQueryKey = (params?: MaybeRef<GetUsernamesByPrefixQueryParams>) =>
  [{ url: '/data/users_by_prefix' }, ...(params ? [params] : [])] as const

export type GetUsernamesByPrefixQueryKey = ReturnType<typeof getUsernamesByPrefixQueryKey>

/**
 * @description Retrieve a list of usernames that match the given prefix
 * @summary Get usernames by prefix
 * {@link /data/users_by_prefix}
 */
async function getUsernamesByPrefix(params?: GetUsernamesByPrefixQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetUsernamesByPrefixQueryResponse, ResponseErrorConfig<Error>, unknown>({
    method: 'GET',
    url: '/data/users_by_prefix',
    params,
    ...config
  })
  
  return res.data
}

export function getUsernamesByPrefixQueryOptions(params?: MaybeRef<GetUsernamesByPrefixQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getUsernamesByPrefixQueryKey(params)
  
  return queryOptions<GetUsernamesByPrefixQueryResponse, ResponseErrorConfig<Error>, GetUsernamesByPrefixQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getUsernamesByPrefix(unref(params), unref(config))
    }
  })
}

/**
 * @description Retrieve a list of usernames that match the given prefix
 * @summary Get usernames by prefix
 * {@link /data/users_by_prefix}
 */
export function useGetUsernamesByPrefix<
  TData = GetUsernamesByPrefixQueryResponse,
  TQueryData = GetUsernamesByPrefixQueryResponse,
  TQueryKey extends QueryKey = GetUsernamesByPrefixQueryKey,
>(
  params?: MaybeRef<GetUsernamesByPrefixQueryParams>,
  options: {
    query?: Partial<QueryObserverOptions<GetUsernamesByPrefixQueryResponse, ResponseErrorConfig<Error>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getUsernamesByPrefixQueryKey(params)

  const query = useQuery({
    ...(getUsernamesByPrefixQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<Error>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}