import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { GetUsernamesByPrefixQueryResponse, GetUsernamesByPrefixQueryParams } from '../types/GetUsernamesByPrefix.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const getUsernamesByPrefixInfiniteQueryKey = (params?: MaybeRef<GetUsernamesByPrefixQueryParams>) =>
  [{ url: '/data/users_by_prefix' }, ...(params ? [params] : [])] as const

export type GetUsernamesByPrefixInfiniteQueryKey = ReturnType<typeof getUsernamesByPrefixInfiniteQueryKey>

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

export function getUsernamesByPrefixInfiniteQueryOptions(params?: MaybeRef<GetUsernamesByPrefixQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getUsernamesByPrefixInfiniteQueryKey(params)
  
  return infiniteQueryOptions<GetUsernamesByPrefixQueryResponse, ResponseErrorConfig<Error>, GetUsernamesByPrefixQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      if (params) {
        params['next_page'] = pageParam as unknown as GetUsernamesByPrefixQueryParams['next_page']
      }
      
      return getUsernamesByPrefix(params, config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @description Retrieve a list of usernames that match the given prefix
 * @summary Get usernames by prefix
 * {@link /data/users_by_prefix}
 */
export function useGetUsernamesByPrefixInfinite<
  TData = InfiniteData<GetUsernamesByPrefixQueryResponse>,
  TQueryData = GetUsernamesByPrefixQueryResponse,
  TQueryKey extends QueryKey = GetUsernamesByPrefixInfiniteQueryKey,
>(
  params?: MaybeRef<GetUsernamesByPrefixQueryParams>,
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetUsernamesByPrefixQueryResponse, ResponseErrorConfig<Error>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getUsernamesByPrefixInfiniteQueryKey(params)

  const query = useInfiniteQuery({
    ...(getUsernamesByPrefixInfiniteQueryOptions(params, config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, ResponseErrorConfig<Error>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}