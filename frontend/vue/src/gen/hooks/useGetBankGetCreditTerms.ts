import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import { unref } from 'vue'
import type { GetBankGetCreditTermsQueryResponse, GetBankGetCreditTermsQueryParams, GetBankGetCreditTerms500 } from '../types/GetBankGetCreditTerms.ts'

export const getBankGetCreditTermsQueryKey = (params?: MaybeRef<GetBankGetCreditTermsQueryParams>) =>
  [{ url: '/bank/get_credit_terms' }, ...(params ? [params] : [])] as const

export type GetBankGetCreditTermsQueryKey = ReturnType<typeof getBankGetCreditTermsQueryKey>

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

export function getBankGetCreditTermsQueryOptions(params?: MaybeRef<GetBankGetCreditTermsQueryParams>, config: Partial<RequestConfig> = {}) {
  const queryKey = getBankGetCreditTermsQueryKey(params)
  
  return queryOptions<GetBankGetCreditTermsQueryResponse, GetBankGetCreditTerms500, GetBankGetCreditTermsQueryResponse, typeof queryKey>({
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
    query?: Partial<QueryObserverOptions<GetBankGetCreditTermsQueryResponse, GetBankGetCreditTerms500, TData, TQueryData, TQueryKey>>;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBankGetCreditTermsQueryKey(params)

  const query = useQuery({
    ...(getBankGetCreditTermsQueryOptions(params, config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, GetBankGetCreditTerms500> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}