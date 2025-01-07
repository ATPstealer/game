import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { GetEquipmentMyQueryResponse, GetEquipmentMyQueryParams, GetEquipmentMy401, GetEquipmentMy500 } from '../types/GetEquipmentMy.ts'

export const getEquipmentMyInfiniteQueryKey = (params?: MaybeRef<GetEquipmentMyQueryParams>) => [{ url: '/equipment/my' }, ...(params ? [params] : [])] as const

export type GetEquipmentMyInfiniteQueryKey = ReturnType<typeof getEquipmentMyInfiniteQueryKey>

/**
 * @summary Return user's equipment
 * {@link /equipment/my}
 */
async function getEquipmentMy(params?: GetEquipmentMyQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetEquipmentMyQueryResponse, GetEquipmentMy401 | GetEquipmentMy500, unknown>({
    method: 'GET',
    url: '/equipment/my',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    params,
    ...config
  })
  
  return res.data
}

export function getEquipmentMyInfiniteQueryOptions(params?: MaybeRef<GetEquipmentMyQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getEquipmentMyInfiniteQueryKey(params)
  
  return infiniteQueryOptions<GetEquipmentMyQueryResponse, GetEquipmentMy401 | GetEquipmentMy500, GetEquipmentMyQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      if (params) {
        params['next_page'] = pageParam as unknown as GetEquipmentMyQueryParams['next_page']
      }
      
      return getEquipmentMy(params, config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Return user's equipment
 * {@link /equipment/my}
 */
export function useGetEquipmentMyInfinite<
  TData = InfiniteData<GetEquipmentMyQueryResponse>,
  TQueryData = GetEquipmentMyQueryResponse,
  TQueryKey extends QueryKey = GetEquipmentMyInfiniteQueryKey,
>(
  params?: MaybeRef<GetEquipmentMyQueryParams>,
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetEquipmentMyQueryResponse, GetEquipmentMy401 | GetEquipmentMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getEquipmentMyInfiniteQueryKey(params)

  const query = useInfiniteQuery({
    ...(getEquipmentMyInfiniteQueryOptions(params, config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetEquipmentMy401 | GetEquipmentMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}