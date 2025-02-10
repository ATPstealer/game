import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetUserDataQueryResponse, GetUserData401 } from '../types/GetUserData.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const getUserDataInfiniteQueryKey = () => [{ url: '/user/data' }] as const

export type GetUserDataInfiniteQueryKey = ReturnType<typeof getUserDataInfiniteQueryKey>

/**
 * @summary Get user data
 * {@link /user/data}
 */
async function getUserData(config: Partial<RequestConfig> = {}) {
  const res = await client<GetUserDataQueryResponse, ResponseErrorConfig<GetUserData401>, unknown>({ method: 'GET', url: '/user/data', ...config })
  
  return res.data
}

export function getUserDataInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getUserDataInfiniteQueryKey()
  
  return infiniteQueryOptions<GetUserDataQueryResponse, ResponseErrorConfig<GetUserData401>, GetUserDataQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getUserData(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Get user data
 * {@link /user/data}
 */
export function useGetUserDataInfinite<
  TData = InfiniteData<GetUserDataQueryResponse>,
  TQueryData = GetUserDataQueryResponse,
  TQueryKey extends QueryKey = GetUserDataInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetUserDataQueryResponse, ResponseErrorConfig<GetUserData401>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getUserDataInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getUserDataInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, ResponseErrorConfig<GetUserData401>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}