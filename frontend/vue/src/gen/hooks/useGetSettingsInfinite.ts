import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetSettingsQueryResponse, GetSettings500 } from '../types/GetSettings.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const getSettingsInfiniteQueryKey = () => [{ url: '/settings' }] as const

export type GetSettingsInfiniteQueryKey = ReturnType<typeof getSettingsInfiniteQueryKey>

/**
 * @description X Y dimension, Interest rate, etc
 * @summary Get General Game Settings
 * {@link /settings}
 */
async function getSettings(config: Partial<RequestConfig> = {}) {
  const res = await client<GetSettingsQueryResponse, ResponseErrorConfig<GetSettings500>, unknown>({ method: 'GET', url: '/settings', ...config })
  
  return res.data
}

export function getSettingsInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getSettingsInfiniteQueryKey()
  
  return infiniteQueryOptions<GetSettingsQueryResponse, ResponseErrorConfig<GetSettings500>, GetSettingsQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getSettings(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @description X Y dimension, Interest rate, etc
 * @summary Get General Game Settings
 * {@link /settings}
 */
export function useGetSettingsInfinite<
  TData = InfiniteData<GetSettingsQueryResponse>,
  TQueryData = GetSettingsQueryResponse,
  TQueryKey extends QueryKey = GetSettingsInfiniteQueryKey,
>(
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetSettingsQueryResponse, ResponseErrorConfig<GetSettings500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getSettingsInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getSettingsInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, ResponseErrorConfig<GetSettings500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}