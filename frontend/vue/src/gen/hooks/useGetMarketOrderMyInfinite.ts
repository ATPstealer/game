import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetMarketOrderMyQueryResponse, GetMarketOrderMy401, GetMarketOrderMy500 } from '../types/GetMarketOrderMy.ts'

export const getMarketOrderMyInfiniteQueryKey = () => [{ url: '/market/order/my' }] as const

export type GetMarketOrderMyInfiniteQueryKey = ReturnType<typeof getMarketOrderMyInfiniteQueryKey>

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

export function getMarketOrderMyInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getMarketOrderMyInfiniteQueryKey()
  
  return infiniteQueryOptions<GetMarketOrderMyQueryResponse, GetMarketOrderMy401 | GetMarketOrderMy500, GetMarketOrderMyQueryResponse, typeof queryKey, number>(
    {
      queryKey,
      queryFn: async ({ signal, pageParam }) => {
        config.signal = signal

        return getMarketOrderMy(config)
      },
      initialPageParam: 0,
      getNextPageParam: (lastPage) => lastPage['nextCursor'],
      getPreviousPageParam: (firstPage) => firstPage['nextCursor']
    }
  )
}

/**
 * @summary Get my orders
 * {@link /market/order/my}
 */
export function useGetMarketOrderMyInfinite<
  TData = InfiniteData<GetMarketOrderMyQueryResponse>,
  TQueryData = GetMarketOrderMyQueryResponse,
  TQueryKey extends QueryKey = GetMarketOrderMyInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetMarketOrderMyQueryResponse, GetMarketOrderMy401 | GetMarketOrderMy500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getMarketOrderMyInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getMarketOrderMyInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetMarketOrderMy401 | GetMarketOrderMy500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}