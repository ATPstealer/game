import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetMapCellOwnersQueryResponse, GetMapCellOwnersQueryParams, GetMapCellOwners500 } from '../types/GetMapCellOwners.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getMapCellOwnersQueryKey = (params: MaybeRef<GetMapCellOwnersQueryParams>) => [{ url: '/map/cell_owners' }, ...(params ? [params] : [])] as const

export type GetMapCellOwnersQueryKey = ReturnType<typeof getMapCellOwnersQueryKey>

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

export function getMapCellOwnersQueryOptions(params: MaybeRef<GetMapCellOwnersQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getMapCellOwnersQueryKey(params)
  
  return queryOptions<GetMapCellOwnersQueryResponse, ResponseErrorConfig<GetMapCellOwners500>, GetMapCellOwnersQueryResponse, typeof queryKey>({
    enabled: !!params,
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getMapCellOwners(unref(params), unref(config))
    }
  })
}

/**
 * @summary Get the landlords in cell
 * {@link /map/cell_owners}
 */
export function useGetMapCellOwners<
  TData = GetMapCellOwnersQueryResponse,
  TQueryData = GetMapCellOwnersQueryResponse,
  TQueryKey extends QueryKey = GetMapCellOwnersQueryKey,
>(
  params: MaybeRef<GetMapCellOwnersQueryParams>,
  options: {
    query?: Partial<QueryObserverOptions<GetMapCellOwnersQueryResponse, ResponseErrorConfig<GetMapCellOwners500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMapCellOwnersQueryKey(params)

  const query = useQuery({
    ...(getMapCellOwnersQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetMapCellOwners500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}