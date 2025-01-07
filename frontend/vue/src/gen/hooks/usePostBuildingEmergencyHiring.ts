import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBuildingEmergencyHiringMutationRequest,
  PostBuildingEmergencyHiringMutationResponse,
  PostBuildingEmergencyHiring401,
  PostBuildingEmergencyHiring500
} from '../types/PostBuildingEmergencyHiring.ts'

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
    PostBuildingEmergencyHiring401 | PostBuildingEmergencyHiring500,
    PostBuildingEmergencyHiringMutationRequest
  >({ method: 'POST', url: '/building/emergency_hiring', baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2', data, ...config })
  
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
      PostBuildingEmergencyHiring401 | PostBuildingEmergencyHiring500,
      { data: MaybeRef<PostBuildingEmergencyHiringMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBuildingEmergencyHiringMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBuildingEmergencyHiringMutationKey()

  return useMutation<
    PostBuildingEmergencyHiringMutationResponse,
    PostBuildingEmergencyHiring401 | PostBuildingEmergencyHiring500,
    { data: PostBuildingEmergencyHiringMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postBuildingEmergencyHiring(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}