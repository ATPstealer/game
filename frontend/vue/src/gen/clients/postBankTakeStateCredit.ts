import client from '@kubb/plugin-client/clients/axios'
import type { RequestConfig } from '@kubb/plugin-client/clients/axios'
import type {
  PostBankTakeStateCreditMutationRequest,
  PostBankTakeStateCreditMutationResponse,
  PostBankTakeStateCredit401,
  PostBankTakeStateCredit500
} from '../types/PostBankTakeStateCredit.ts'

/**
 * @description Get credit from state. Payload example {"buildingId":"670fd64c211de59e1bb8a314", "Amount": 5000}
 * @summary Take state credit
 * {@link /bank/take_state_credit}
 */
export async function postBankTakeStateCredit(
  data: PostBankTakeStateCreditMutationRequest,
  config: Partial<RequestConfig<PostBankTakeStateCreditMutationRequest>> = {}
) {
  const res = await client<
    PostBankTakeStateCreditMutationResponse,
    PostBankTakeStateCredit401 | PostBankTakeStateCredit500,
    PostBankTakeStateCreditMutationRequest
  >({ method: 'POST', url: '/bank/take_state_credit', data, ...config })
  
  return res
}