import client from '@kubb/plugin-client/dist/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/dist/clients/axios'
import type {
  PostBankTakeCreditMutationRequest,
  PostBankTakeCreditMutationResponse,
  PostBankTakeCredit401,
  PostBankTakeCredit500
} from '../types/PostBankTakeCredit.ts'

/**
 * @description Get credit in bank. Payload example {"buildingId":"670fd64c211de59e1bb8a314", "Amount":50, "Rate": 0.5, "Rating": -1000000}
 * @summary Take credit
 * {@link /bank/take_credit}
 */
export async function postBankTakeCredit(data: PostBankTakeCreditMutationRequest, config: Partial<RequestConfig<PostBankTakeCreditMutationRequest>> = {}) {
  const res = await client<PostBankTakeCreditMutationResponse, PostBankTakeCredit401 | PostBankTakeCredit500, PostBankTakeCreditMutationRequest>({
    method: 'POST',
    url: '/bank/take_credit',
    data,
    ...config
  })
  
  return res.data
}