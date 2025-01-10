import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetResourceMyQueryResponse, GetResourceMyQueryParams, GetResourceMy401, GetResourceMy500 } from '../types/GetResourceMy.ts'
import type { RequestConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getResourceMyQueryKey = (params?: MaybeRef<GetResourceMyQueryParams>) => [{ url: '/resource/my' }, ...(params ? [params] : [])] as const

export type GetResourceMyQueryKey = ReturnType<typeof getResourceMyQueryKey>

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

export function getResourceMyQueryOptions(params?: MaybeRef<GetResourceMyQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getResourceMyQueryKey(params)
  
  return queryOptions<GetResourceMyQueryResponse, GetResourceMy401 | GetResourceMy500, GetResourceMyQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getResourceMy(unref(params), unref(config))
    }
  })
}

/**
 * @summary Get user's resources
 * {@link /resource/my}
 */
export function useGetResourceMy<
  TData = GetResourceMyQueryResponse,
  TQueryData = GetResourceMyQueryResponse,
  TQueryKey extends QueryKey = GetResourceMyQueryKey,
>(
  params?: MaybeRef<GetResourceMyQueryParams>,
  options: {
    query?: Partial<QueryObserverOptions<GetResourceMyQueryResponse, GetResourceMy401 | GetResourceMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getResourceMyQueryKey(params)

  const query = useQuery({
    ...(getResourceMyQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetResourceMy401 | GetResourceMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}