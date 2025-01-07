import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { GetResourceMyQueryResponse, GetResourceMyQueryParams, GetResourceMy401, GetResourceMy500 } from '../types/GetResourceMy.ts'

export const getResourceMyInfiniteQueryKey = (params?: MaybeRef<GetResourceMyQueryParams>) => [{ url: '/resource/my' }, ...(params ? [params] : [])] as const

export type GetResourceMyInfiniteQueryKey = ReturnType<typeof getResourceMyInfiniteQueryKey>

/**
 * @summary Get user's resources
 * {@link /resource/my}
 */
async function getResourceMy(params?: GetResourceMyQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetResourceMyQueryResponse, GetResourceMy401 | GetResourceMy500, unknown>({
    method: 'GET',
    url: '/resource/my',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    params,
    ...config
  })
  
  return res.data
}

export function getResourceMyInfiniteQueryOptions(params?: MaybeRef<GetResourceMyQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getResourceMyInfiniteQueryKey(params)
  
  return infiniteQueryOptions<GetResourceMyQueryResponse, GetResourceMy401 | GetResourceMy500, GetResourceMyQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      if (params) {
        params['next_page'] = pageParam as unknown as GetResourceMyQueryParams['next_page']
      }
      
      return getResourceMy(params, config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Get user's resources
 * {@link /resource/my}
 */
export function useGetResourceMyInfinite<
  TData = InfiniteData<GetResourceMyQueryResponse>,
  TQueryData = GetResourceMyQueryResponse,
  TQueryKey extends QueryKey = GetResourceMyInfiniteQueryKey,
>(
  params?: MaybeRef<GetResourceMyQueryParams>,
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetResourceMyQueryResponse, GetResourceMy401 | GetResourceMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getResourceMyInfiniteQueryKey(params)

  const query = useInfiniteQuery({
    ...(getResourceMyInfiniteQueryOptions(params, config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetResourceMy401 | GetResourceMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}