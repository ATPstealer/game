import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { GetMapCellOwnersQueryResponse, GetMapCellOwnersQueryParams, GetMapCellOwners500 } from '../types/GetMapCellOwners.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const getMapCellOwnersInfiniteQueryKey = (params: MaybeRef<GetMapCellOwnersQueryParams>) =>
  [{ url: '/map/cell_owners' }, ...(params ? [params] : [])] as const

export type GetMapCellOwnersInfiniteQueryKey = ReturnType<typeof getMapCellOwnersInfiniteQueryKey>

/**
 * @summary Get the landlords in cell
 * {@link /map/cell_owners}
 */
async function getMapCellOwners(params: GetMapCellOwnersQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetMapCellOwnersQueryResponse, ResponseErrorConfig<GetMapCellOwners500>, unknown>({
    method: 'GET',
    url: '/map/cell_owners',
    params,
    ...config
  })
  
  return res.data
}

export function getMapCellOwnersInfiniteQueryOptions(params: MaybeRef<GetMapCellOwnersQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getMapCellOwnersInfiniteQueryKey(params)
  
  return infiniteQueryOptions<GetMapCellOwnersQueryResponse, ResponseErrorConfig<GetMapCellOwners500>, GetMapCellOwnersQueryResponse, typeof queryKey, number>({
    enabled: !!params,
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      if (params) {
        params['next_page'] = pageParam as unknown as GetMapCellOwnersQueryParams['next_page']
      }
      
      return getMapCellOwners(params, config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Get the landlords in cell
 * {@link /map/cell_owners}
 */
export function useGetMapCellOwnersInfinite<
  TData = InfiniteData<GetMapCellOwnersQueryResponse>,
  TQueryData = GetMapCellOwnersQueryResponse,
  TQueryKey extends QueryKey = GetMapCellOwnersInfiniteQueryKey,
>(
  params: MaybeRef<GetMapCellOwnersQueryParams>,
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetMapCellOwnersQueryResponse, ResponseErrorConfig<GetMapCellOwners500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMapCellOwnersInfiniteQueryKey(params)

  const query = useInfiniteQuery({
    ...(getMapCellOwnersInfiniteQueryOptions(params, config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, ResponseErrorConfig<GetMapCellOwners500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}