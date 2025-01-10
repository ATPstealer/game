import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetEquipmentTypesQueryResponse, GetEquipmentTypes500 } from '../types/GetEquipmentTypes.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const getEquipmentTypesInfiniteQueryKey = () => [{ url: '/equipment/types' }] as const

export type GetEquipmentTypesInfiniteQueryKey = ReturnType<typeof getEquipmentTypesInfiniteQueryKey>

/**
 * @summary Get all equipment types
 * {@link /equipment/types}
 */
async function getEquipmentTypes(config: Partial<RequestConfig> = {}) {
  const res = await client<GetEquipmentTypesQueryResponse, GetEquipmentTypes500, unknown>({
    method: 'GET',
    url: '/equipment/types',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getEquipmentTypesInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getEquipmentTypesInfiniteQueryKey()
  
  return infiniteQueryOptions<GetEquipmentTypesQueryResponse, GetEquipmentTypes500, GetEquipmentTypesQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getEquipmentTypes(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @summary Get all equipment types
 * {@link /equipment/types}
 */
export function useGetEquipmentTypesInfinite<
  TData = InfiniteData<GetEquipmentTypesQueryResponse>,
  TQueryData = GetEquipmentTypesQueryResponse,
  TQueryKey extends QueryKey = GetEquipmentTypesInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetEquipmentTypesQueryResponse, GetEquipmentTypes500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getEquipmentTypesInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getEquipmentTypesInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetEquipmentTypes500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}