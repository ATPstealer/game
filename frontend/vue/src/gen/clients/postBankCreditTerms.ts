import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type {
  PostBankCreditTermsMutationRequest,
  PostBankCreditTermsMutationResponse,
  PostBankCreditTerms401,
  PostBankCreditTerms500
} from '../types/PostBankCreditTerms.ts'

/**
 * @description Limit > 0, Rate > 0. For change limit send payload: {"Rate": sameAsExisting, "Rating": sameAsExisting, "Adding": true}
 * @summary Add / Change / Delete credit terms in bank contracts
 * {@link /bank/credit_terms}
 */
export async function postBankCreditTerms(data: PostBankCreditTermsMutationRequest, config: Partial<RequestConfig<PostBankCreditTermsMutationRequest>> = {}) {
  const res = await client<PostBankCreditTermsMutationResponse, PostBankCreditTerms401 | PostBankCreditTerms500, PostBankCreditTermsMutationRequest>({
    method: 'POST',
    url: '/bank/credit_terms',
    data,
    ...config
  })
  
  return res
}