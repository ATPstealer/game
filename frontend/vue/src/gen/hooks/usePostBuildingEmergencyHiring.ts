import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBuildingEmergencyHiringMutationRequest,
  PostBuildingEmergencyHiringMutationResponse,
  PostBuildingEmergencyHiring401,
  PostBuildingEmergencyHiring500
} from '../types/PostBuildingEmergencyHiring.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const postBuildingEmergencyHiringMutationKey = () => [{ url: '/building/emergency_hiring' }] as const

export type PostBuildingEmergencyHiringMutationKey = ReturnType<typeof postBuildingEmergencyHiringMutationKey>

/**
 * @summary Expensive fast hiring
 * {@link /building/emergency_hiring}
 */
async function postBuildingEmergencyHiring(
  data: PostBuildingEmergencyHiringMutationRequest,
  config: Partial<RequestConfig<PostBuildingEmergencyHiringMutationRequest>> = {}
) {
  const res = await client<
    PostBuildingEmergencyHiringMutationResponse,
    ResponseErrorConfig<PostBuildingEmergencyHiring401 | PostBuildingEmergencyHiring500>,
    PostBuildingEmergencyHiringMutationRequest
  >({ method: 'POST', url: '/building/emergency_hiring', data, ...config })
  
  return res.data
}

/**
 * @summary Expensive fast hiring
 * {@link /building/emergency_hiring}
 */
export function usePostBuildingEmergencyHiring(
  options: {
    mutation?: MutationObserverOptions<
      PostBuildingEmergencyHiringMutationResponse,
      ResponseErrorConfig<PostBuildingEmergencyHiring401 | PostBuildingEmergencyHiring500>,
      { data: MaybeRef<PostBuildingEmergencyHiringMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBuildingEmergencyHiringMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBuildingEmergencyHiringMutationKey()

  return useMutation<
    PostBuildingEmergencyHiringMutationResponse,
    ResponseErrorConfig<PostBuildingEmergencyHiring401 | PostBuildingEmergencyHiring500>,
    { data: PostBuildingEmergencyHiringMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postBuildingEmergencyHiring(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}