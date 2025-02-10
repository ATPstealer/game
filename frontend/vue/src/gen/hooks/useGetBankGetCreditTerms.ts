import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetBankGetCreditTermsQueryResponse, GetBankGetCreditTermsQueryParams, GetBankGetCreditTerms500 } from '../types/GetBankGetCreditTerms.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getBankGetCreditTermsQueryKey = (params?: MaybeRef<GetBankGetCreditTermsQueryParams>) =>
  [{ url: '/bank/get_credit_terms' }, ...(params ? [params] : [])] as const

export type GetBankGetCreditTermsQueryKey = ReturnType<typeof getBankGetCreditTermsQueryKey>

/**
 * @description If defined return. Credit term where limit >= in param, rate <= in param, rating <= in param.
 * @summary Return credit terms
 * {@link /bank/get_credit_terms}
 */
async function getBankGetCreditTerms(params?: GetBankGetCreditTermsQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<GetBankGetCreditTermsQueryResponse, ResponseErrorConfig<GetBankGetCreditTerms500>, unknown>({
    method: 'GET',
    url: '/bank/get_credit_terms',
    params,
    ...config
  })
  
  return res.data
}

export function getBankGetCreditTermsQueryOptions(params?: MaybeRef<GetBankGetCreditTermsQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getBankGetCreditTermsQueryKey(params)
  
  return queryOptions<GetBankGetCreditTermsQueryResponse, ResponseErrorConfig<GetBankGetCreditTerms500>, GetBankGetCreditTermsQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getBankGetCreditTerms(unref(params), unref(config))
    }
  })
}

/**
 * @description If defined return. Credit term where limit >= in param, rate <= in param, rating <= in param.
 * @summary Return credit terms
 * {@link /bank/get_credit_terms}
 */
export function useGetBankGetCreditTerms<
  TData = GetBankGetCreditTermsQueryResponse,
  TQueryData = GetBankGetCreditTermsQueryResponse,
  TQueryKey extends QueryKey = GetBankGetCreditTermsQueryKey,
>(
  params?: MaybeRef<GetBankGetCreditTermsQueryParams>,
  options: {
    query?: Partial<QueryObserverOptions<GetBankGetCreditTermsQueryResponse, ResponseErrorConfig<GetBankGetCreditTerms500>, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBankGetCreditTermsQueryKey(params)

  const query = useQuery({
    ...(getBankGetCreditTermsQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetBankGetCreditTerms500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}