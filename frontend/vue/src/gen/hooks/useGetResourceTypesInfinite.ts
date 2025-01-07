import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetResourceTypesQueryResponse, GetResourceTypes500 } from '../types/GetResourceTypes.ts'

export const getResourceTypesInfiniteQueryKey = () => [{ url: '/resource/types' }] as const

export type GetResourceTypesInfiniteQueryKey = ReturnType<typeof getResourceTypesInfiniteQueryKey>

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

export function getResourceTypesInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getResourceTypesInfiniteQueryKey()
  
  return infiniteQueryOptions<GetResourceTypesQueryResponse, GetResourceTypes500, GetResourceTypesQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getResourceTypes(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Return all resource types from database
 * {@link /resource/types}
 */
export function useGetResourceTypesInfinite<
  TData = InfiniteData<GetResourceTypesQueryResponse>,
  TQueryData = GetResourceTypesQueryResponse,
  TQueryKey extends QueryKey = GetResourceTypesInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetResourceTypesQueryResponse, GetResourceTypes500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getResourceTypesInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getResourceTypesInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetResourceTypes500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}