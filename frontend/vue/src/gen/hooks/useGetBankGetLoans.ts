import type { QueryKey, QueryObserverOptions, UseQueryReturnType } from '@tanstack/vue-query'
import { queryOptions, useQuery } from '@tanstack/vue-query'
import { unref } from 'vue'
import type { GetBankGetLoansQueryResponse, GetBankGetLoans401, GetBankGetLoans500 } from '../types/GetBankGetLoans.ts'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'
import client from '@/api/customClientAxios'

export const getBankGetLoansQueryKey = () => [{ url: '/bank/get_loans' }] as const

export type GetBankGetLoansQueryKey = ReturnType<typeof getBankGetLoansQueryKey>

/**
 * @description Return all loans connected with userId
 * @summary Get Users Loans
 * {@link /bank/get_loans}
 */
async function getBankGetLoans(config: Partial<RequestConfig> = {}) {
  const res = await client<GetBankGetLoansQueryResponse, ResponseErrorConfig<GetBankGetLoans401 | GetBankGetLoans500>, unknown>({
    method: 'GET',
    url: '/bank/get_loans',
    ...config
  })
  
  return res.data
}

export function getBankGetLoansQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getBankGetLoansQueryKey()
  
  return queryOptions<
    GetBankGetLoansQueryResponse,
    ResponseErrorConfig<GetBankGetLoans401 | GetBankGetLoans500>,
    GetBankGetLoansQueryResponse,
    typeof queryKey
  >({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      
      return getBankGetLoans(unref(config))
    }
  })
}

/**
 * @description Return all loans connected with userId
 * @summary Get Users Loans
 * {@link /bank/get_loans}
 */
export function useGetBankGetLoans<
  TData = GetBankGetLoansQueryResponse,
  TQueryData = GetBankGetLoansQueryResponse,
  TQueryKey extends QueryKey = GetBankGetLoansQueryKey,
>(
  options: {
    query?: Partial<
      QueryObserverOptions<GetBankGetLoansQueryResponse, ResponseErrorConfig<GetBankGetLoans401 | GetBankGetLoans500>, TData, TQueryData, TQueryKey>
    >;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBankGetLoansQueryKey()

  const query = useQuery({
    ...(getBankGetLoansQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>)
  }) as UseQueryReturnType<TData, ResponseErrorConfig<GetBankGetLoans401 | GetBankGetLoans500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}