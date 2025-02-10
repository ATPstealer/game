import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBankTakeStateCreditMutationRequest,
  PostBankTakeStateCreditMutationResponse,
  PostBankTakeStateCredit401,
  PostBankTakeStateCredit500
} from '../types/PostBankTakeStateCredit.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const postBankTakeStateCreditMutationKey = () => [{ url: '/bank/take_state_credit' }] as const

export type PostBankTakeStateCreditMutationKey = ReturnType<typeof postBankTakeStateCreditMutationKey>

/**
 * @description Get credit from state. Payload example {"buildingId":"670fd64c211de59e1bb8a314", "Amount": 5000}
 * @summary Take state credit
 * {@link /bank/take_state_credit}
 */
async function postBankTakeStateCredit(
  data: PostBankTakeStateCreditMutationRequest,
  config: Partial<RequestConfig<PostBankTakeStateCreditMutationRequest>> = {}
) {
  const res = await client<
    PostBankTakeStateCreditMutationResponse,
    ResponseErrorConfig<PostBankTakeStateCredit401 | PostBankTakeStateCredit500>,
    PostBankTakeStateCreditMutationRequest
  >({ method: 'POST', url: '/bank/take_state_credit', data, ...config })
  
  return res.data
}

/**
 * @description Get credit from state. Payload example {"buildingId":"670fd64c211de59e1bb8a314", "Amount": 5000}
 * @summary Take state credit
 * {@link /bank/take_state_credit}
 */
export function usePostBankTakeStateCredit(
  options: {
    mutation?: MutationObserverOptions<
      PostBankTakeStateCreditMutationResponse,
      ResponseErrorConfig<PostBankTakeStateCredit401 | PostBankTakeStateCredit500>,
      { data: MaybeRef<PostBankTakeStateCreditMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBankTakeStateCreditMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBankTakeStateCreditMutationKey()

  return useMutation<
    PostBankTakeStateCreditMutationResponse,
    ResponseErrorConfig<PostBankTakeStateCredit401 | PostBankTakeStateCredit500>,
    { data: PostBankTakeStateCreditMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postBankTakeStateCredit(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}