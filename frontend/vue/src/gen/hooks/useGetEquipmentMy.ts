import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetEquipmentMyQueryResponse, GetEquipmentMyQueryParams, GetEquipmentMy401, GetEquipmentMy500 } from '../types/GetEquipmentMy.ts'
import type { RequestConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getEquipmentMyQueryKey = (params?: MaybeRef<GetEquipmentMyQueryParams>) => [{ url: '/equipment/my' }, ...(params ? [params] : [])] as const

export type GetEquipmentMyQueryKey = ReturnType<typeof getEquipmentMyQueryKey>

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

export function getEquipmentMyQueryOptions(params?: MaybeRef<GetEquipmentMyQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getEquipmentMyQueryKey(params)
  
  return queryOptions<GetEquipmentMyQueryResponse, GetEquipmentMy401 | GetEquipmentMy500, GetEquipmentMyQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getEquipmentMy(unref(params), unref(config))
    }
  })
}

/**
 * @summary Return user's equipment
 * {@link /equipment/my}
 */
export function useGetEquipmentMy<
  TData = GetEquipmentMyQueryResponse,
  TQueryData = GetEquipmentMyQueryResponse,
  TQueryKey extends QueryKey = GetEquipmentMyQueryKey,
>(
  params?: MaybeRef<GetEquipmentMyQueryParams>,
  options: {
    query?: Partial<QueryObserverOptions<GetEquipmentMyQueryResponse, GetEquipmentMy401 | GetEquipmentMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getEquipmentMyQueryKey(params)

  const query = useQuery({
    ...(getEquipmentMyQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetEquipmentMy401 | GetEquipmentMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}