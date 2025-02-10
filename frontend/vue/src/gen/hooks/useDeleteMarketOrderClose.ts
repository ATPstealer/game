import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  DeleteMarketOrderCloseMutationResponse,
  DeleteMarketOrderCloseQueryParams,
  DeleteMarketOrderClose401,
  DeleteMarketOrderClose500
} from '../types/DeleteMarketOrderClose.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const deleteMarketOrderCloseMutationKey = () => [{ url: '/market/order/close' }] as const

export type DeleteMarketOrderCloseMutationKey = ReturnType<typeof deleteMarketOrderCloseMutationKey>

/**
 * @summary Close user's order
 * {@link /market/order/close}
 */
async function deleteMarketOrderClose(params: DeleteMarketOrderCloseQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<DeleteMarketOrderCloseMutationResponse, ResponseErrorConfig<DeleteMarketOrderClose401 | DeleteMarketOrderClose500>, unknown>({
    method: 'DELETE',
    url: '/market/order/close',
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
      ResponseErrorConfig<DeleteMarketOrderClose401 | DeleteMarketOrderClose500>,
      { params: MaybeRef<DeleteMarketOrderCloseQueryParams> }
    >;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? deleteMarketOrderCloseMutationKey()

  return useMutation<
    DeleteMarketOrderCloseMutationResponse,
    ResponseErrorConfig<DeleteMarketOrderClose401 | DeleteMarketOrderClose500>,
    { params: DeleteMarketOrderCloseQueryParams }
  >({
    mutationFn: async ({ params }) => {
      return deleteMarketOrderClose(params, config)
    },
    mutationKey,
    ...mutationOptions
  })
}