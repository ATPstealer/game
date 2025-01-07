import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBankCreditTermsMutationRequest,
  PostBankCreditTermsMutationResponse,
  PostBankCreditTerms401,
  PostBankCreditTerms500
} from '../types/PostBankCreditTerms.ts'

export const postBankCreditTermsMutationKey = () => [{ url: '/bank/credit_terms' }] as const

export type PostBankCreditTermsMutationKey = ReturnType<typeof postBankCreditTermsMutationKey>

/**
 * @description Limit > 0, Rate > 0. For change limit send payload: {"Rate": sameAsExisting, "Rating": sameAsExisting, "Adding": true}
 * @summary Add / Change / Delete credit terms in bank contracts
 * {@link /bank/credit_terms}
 */
async function postBankCreditTerms(data: PostBankCreditTermsMutationRequest, config: Partial<RequestConfig<PostBankCreditTermsMutationRequest>> = {}) {
  const res = await client<PostBankCreditTermsMutationResponse, PostBankCreditTerms401 | PostBankCreditTerms500, PostBankCreditTermsMutationRequest>({
    method: 'POST',
    url: '/bank/credit_terms',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @description Limit > 0, Rate > 0. For change limit send payload: {"Rate": sameAsExisting, "Rating": sameAsExisting, "Adding": true}
 * @summary Add / Change / Delete credit terms in bank contracts
 * {@link /bank/credit_terms}
 */
export function usePostBankCreditTerms(
  options: {
    mutation?: MutationObserverOptions<
      PostBankCreditTermsMutationResponse,
      PostBankCreditTerms401 | PostBankCreditTerms500,
      { data: MaybeRef<PostBankCreditTermsMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBankCreditTermsMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBankCreditTermsMutationKey()

  return useMutation<PostBankCreditTermsMutationResponse, PostBankCreditTerms401 | PostBankCreditTerms500, { data: PostBankCreditTermsMutationRequest }>({
    mutationFn: async ({ data }) => {
      return postBankCreditTerms(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}