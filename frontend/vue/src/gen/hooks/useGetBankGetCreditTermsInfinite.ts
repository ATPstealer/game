import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { GetBankGetCreditTermsQueryResponse, GetBankGetCreditTermsQueryParams, GetBankGetCreditTerms500 } from '../types/GetBankGetCreditTerms.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const getBankGetCreditTermsInfiniteQueryKey = (params?: MaybeRef<GetBankGetCreditTermsQueryParams>) =>
  [{ url: '/bank/get_credit_terms' }, ...(params ? [params] : [])] as const

export type GetBankGetCreditTermsInfiniteQueryKey = ReturnType<typeof getBankGetCreditTermsInfiniteQueryKey>

/**
 * @description If defined return. Credit term where limit >= in param, rate <= in param, rating <= in param.
 * @summary Return credit terms
 * {@link /bank/get_credit_terms}
 */
async function getBankGetCreditTerms(params?: GetBankGetCreditTermsQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetBankGetCreditTermsQueryResponse, GetBankGetCreditTerms500, unknown>({
    method: 'GET',
    url: '/bank/get_credit_terms',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    params,
    ...config
  })
  
  return res.data
}

export function getBankGetCreditTermsInfiniteQueryOptions(params?: MaybeRef<GetBankGetCreditTermsQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getBankGetCreditTermsInfiniteQueryKey(params)
  
  return infiniteQueryOptions<GetBankGetCreditTermsQueryResponse, GetBankGetCreditTerms500, GetBankGetCreditTermsQueryResponse, typeof queryKey, number>({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      if (params) {
        params['next_page'] = pageParam as unknown as GetBankGetCreditTermsQueryParams['next_page']
      }
      
      return getBankGetCreditTerms(params, config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @description If defined return. Credit term where limit >= in param, rate <= in param, rating <= in param.
 * @summary Return credit terms
 * {@link /bank/get_credit_terms}
 */
export function useGetBankGetCreditTermsInfinite<
  TData = InfiniteData<GetBankGetCreditTermsQueryResponse>,
  TQueryData = GetBankGetCreditTermsQueryResponse,
  TQueryKey extends QueryKey = GetBankGetCreditTermsInfiniteQueryKey,
>(
  params?: MaybeRef<GetBankGetCreditTermsQueryParams>,
  options: {
    query?: Partial<InfiniteQueryObserverOptions<GetBankGetCreditTermsQueryResponse, GetBankGetCreditTerms500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBankGetCreditTermsInfiniteQueryKey(params)

  const query = useInfiniteQuery({
    ...(getBankGetCreditTermsInfiniteQueryOptions(params, config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, GetBankGetCreditTerms500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}