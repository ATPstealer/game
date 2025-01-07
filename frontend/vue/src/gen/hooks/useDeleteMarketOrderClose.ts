import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  DeleteMarketOrderCloseMutationResponse,
  DeleteMarketOrderCloseQueryParams,
  DeleteMarketOrderClose401,
  DeleteMarketOrderClose500
} from '../types/DeleteMarketOrderClose.ts'

export const deleteMarketOrderCloseMutationKey = () => [{ url: '/market/order/close' }] as const

export type DeleteMarketOrderCloseMutationKey = ReturnType<typeof deleteMarketOrderCloseMutationKey>

/**
 * @summary Close user's order
 * {@link /market/order/close}
 */
async function deleteMarketOrderClose(params: DeleteMarketOrderCloseQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<DeleteMarketOrderCloseMutationResponse, DeleteMarketOrderClose401 | DeleteMarketOrderClose500, unknown>({
    method: 'DELETE',
    url: '/market/order/close',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    params,
    ...config
  })
  
  return res.data
}

/**
 * @summary Close user's order
 * {@link /market/order/close}
 */
export function useDeleteMarketOrderClose(
  options: {
    mutation?: MutationObserverOptions<
      DeleteMarketOrderCloseMutationResponse,
      DeleteMarketOrderClose401 | DeleteMarketOrderClose500,
      { params: MaybeRef<DeleteMarketOrderCloseQueryParams> }
    >;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? deleteMarketOrderCloseMutationKey()

  return useMutation<
    DeleteMarketOrderCloseMutationResponse,
    DeleteMarketOrderClose401 | DeleteMarketOrderClose500,
    { params: DeleteMarketOrderCloseQueryParams }
  >({
    mutationFn: async ({ params }) => {
      return deleteMarketOrderClose(params, config)
    },
    mutationKey,
    ...mutationOptions
  })
}