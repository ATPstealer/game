import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetEquipmentTypesQueryResponse, GetEquipmentTypes500 } from '../types/GetEquipmentTypes.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getEquipmentTypesQueryKey = () => [{ url: '/equipment/types' }] as const

export type GetEquipmentTypesQueryKey = ReturnType<typeof getEquipmentTypesQueryKey>

/**
 * @summary Get all equipment types
 * {@link /equipment/types}
 */
async function getEquipmentTypes(config: Partial<RequestConfig> = {}) {
  const res = await client<GetEquipmentTypesQueryResponse, ResponseErrorConfig<GetEquipmentTypes500>, unknown>({
    method: 'GET',
    url: '/equipment/types',
    ...config
  })
  
  return res.data
}

export function getEquipmentTypesQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getEquipmentTypesQueryKey()
  
  return queryOptions<GetEquipmentTypesQueryResponse, ResponseErrorConfig<GetEquipmentTypes500>, GetEquipmentTypesQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getEquipmentTypes(unref(config))
    }
  })
}

/**
 * @summary Get all equipment types
 * {@link /equipment/types}
 */
export function useGetEquipmentTypes<
  TData = GetEquipmentTypesQueryResponse,
  TQueryData = GetEquipmentTypesQueryResponse,
  TQueryKey extends QueryKey = GetEquipmentTypesQueryKey,
>(
  options: {
    query?: Partial<QueryObserverOptions<GetEquipmentTypesQueryResponse, ResponseErrorConfig<GetEquipmentTypes500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getEquipmentTypesQueryKey()

  const query = useQuery({
    ...(getEquipmentTypesQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetEquipmentTypes500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}