import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetMarketOrderMyQueryResponse, GetMarketOrderMy401, GetMarketOrderMy500 } from '../types/GetMarketOrderMy.ts'

export const getMarketOrderMyQueryKey = () => [{ url: '/market/order/my' }] as const

export type GetMarketOrderMyQueryKey = ReturnType<typeof getMarketOrderMyQueryKey>

/**
 * @summary Get my orders
 * {@link /market/order/my}
 */
async function getMarketOrderMy(config: Partial<RequestConfig> = {}) {
  const res = await client<GetMarketOrderMyQueryResponse, GetMarketOrderMy401 | GetMarketOrderMy500, unknown>({
    method: 'GET',
    url: '/market/order/my',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    ...config
  })
  
  return res.data
}

export function getMarketOrderMyQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getMarketOrderMyQueryKey()
  
  return queryOptions<GetMarketOrderMyQueryResponse, GetMarketOrderMy401 | GetMarketOrderMy500, GetMarketOrderMyQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getMarketOrderMy(unref(config))
    }
  })
}

/**
 * @summary Get my orders
 * {@link /market/order/my}
 */
export function useGetMarketOrderMy<
  TData = GetMarketOrderMyQueryResponse,
  TQueryData = GetMarketOrderMyQueryResponse,
  TQueryKey extends QueryKey = GetMarketOrderMyQueryKey,
>(
  options: {
    query?: Partial<QueryObserverOptions<GetMarketOrderMyQueryResponse, GetMarketOrderMy401 | GetMarketOrderMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMarketOrderMyQueryKey()

  const query = useQuery({
    ...(getMarketOrderMyQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetMarketOrderMy401 | GetMarketOrderMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}