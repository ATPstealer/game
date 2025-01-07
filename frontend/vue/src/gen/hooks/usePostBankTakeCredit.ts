import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBankTakeCreditMutationRequest,
  PostBankTakeCreditMutationResponse,
  PostBankTakeCredit401,
  PostBankTakeCredit500
} from '../types/PostBankTakeCredit.ts'

export const postBankTakeCreditMutationKey = () => [{ url: '/bank/take_credit' }] as const

export type PostBankTakeCreditMutationKey = ReturnType<typeof postBankTakeCreditMutationKey>

/**
 * @description Get credit in bank. Payload example {"buildingId":"670fd64c211de59e1bb8a314", "Amount":50, "Rate": 0.5, "Rating": -1000000}
 * @summary Take credit
 * {@link /bank/take_credit}
 */
async function postBankTakeCredit(data: PostBankTakeCreditMutationRequest, config: Partial<RequestConfig<PostBankTakeCreditMutationRequest>> = {}) {
  const res = await client<PostBankTakeCreditMutationResponse, PostBankTakeCredit401 | PostBankTakeCredit500, PostBankTakeCreditMutationRequest>({
    method: 'POST',
    url: '/bank/take_credit',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @description Get credit in bank. Payload example {"buildingId":"670fd64c211de59e1bb8a314", "Amount":50, "Rate": 0.5, "Rating": -1000000}
 * @summary Take credit
 * {@link /bank/take_credit}
 */
export function usePostBankTakeCredit(
  options: {
    mutation?: MutationObserverOptions<
      PostBankTakeCreditMutationResponse,
      PostBankTakeCredit401 | PostBankTakeCredit500,
      { data: MaybeRef<PostBankTakeCreditMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBankTakeCreditMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBankTakeCreditMutationKey()

  return useMutation<PostBankTakeCreditMutationResponse, PostBankTakeCredit401 | PostBankTakeCredit500, { data: PostBankTakeCreditMutationRequest }>({
    mutationFn: async ({ data }) => {
      return postBankTakeCredit(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}