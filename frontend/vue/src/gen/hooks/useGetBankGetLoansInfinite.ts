import type { InfiniteData, QueryKey, InfiniteQueryObserverOptions, UseInfiniteQueryReturnType } from '@tanstack/vue-query'
import { infiniteQueryOptions, useInfiniteQuery } from '@tanstack/vue-query'
import type { GetBankGetLoansQueryResponse, GetBankGetLoans401, GetBankGetLoans500 } from '../types/GetBankGetLoans.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const getBankGetLoansInfiniteQueryKey = () => [{ url: '/bank/get_loans' }] as const

export type GetBankGetLoansInfiniteQueryKey = ReturnType<typeof getBankGetLoansInfiniteQueryKey>

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

export function getBankGetLoansInfiniteQueryOptions(config: Partial<RequestConfig> = {}) {
  const queryKey = getBankGetLoansInfiniteQueryKey()
  
  return infiniteQueryOptions<
    GetBankGetLoansQueryResponse,
    ResponseErrorConfig<GetBankGetLoans401 | GetBankGetLoans500>,
    GetBankGetLoansQueryResponse,
    typeof queryKey,
    number
  >({
    queryKey,
    queryFn: async ({ signal, pageParam }) => {
      config.signal = signal

      return getBankGetLoans(config)
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage) => lastPage['nextCursor'],
    getPreviousPageParam: (firstPage) => firstPage['nextCursor']
  })
}

/**
 * @description Return all loans connected with userId
 * @summary Get Users Loans
 * {@link /bank/get_loans}
 */
export function useGetBankGetLoansInfinite<
  TData = InfiniteData<GetBankGetLoansQueryResponse>,
  TQueryData = GetBankGetLoansQueryResponse,
  TQueryKey extends QueryKey = GetBankGetLoansInfiniteQueryKey,
>(
  options: {
    query?: Partial<
      InfiniteQueryObserverOptions<GetBankGetLoansQueryResponse, ResponseErrorConfig<GetBankGetLoans401 | GetBankGetLoans500>, TData, TQueryData, TQueryKey>
    >;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? getBankGetLoansInfiniteQueryKey()

  const query = useInfiniteQuery({
    ...(getBankGetLoansInfiniteQueryOptions(config) as unknown as InfiniteQueryObserverOptions),
    queryKey: queryKey as QueryKey,
    ...(queryOptions as unknown as Omit<InfiniteQueryObserverOptions, 'queryKey'>)
  }) as UseInfiniteQueryReturnType<TData, ResponseErrorConfig<GetBankGetLoans401 | GetBankGetLoans500>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}