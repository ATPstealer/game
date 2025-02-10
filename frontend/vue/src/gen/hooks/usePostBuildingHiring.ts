import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBuildingHiringMutationRequest,
  PostBuildingHiringMutationResponse,
  PostBuildingHiring401,
  PostBuildingHiring500
} from '../types/PostBuildingHiring.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const postBuildingHiringMutationKey = () => [{ url: '/building/hiring' }] as const

export type PostBuildingHiringMutationKey = ReturnType<typeof postBuildingHiringMutationKey>

/**
 * @summary Set hiring details for a building
 * {@link /building/hiring}
 */
async function postBuildingHiring(data?: PostBuildingHiringMutationRequest, config: Partial<RequestConfig<PostBuildingHiringMutationRequest>> = {}) {
  const res = await client<
    PostBuildingHiringMutationResponse,
    ResponseErrorConfig<PostBuildingHiring401 | PostBuildingHiring500>,
    PostBuildingHiringMutationRequest
  >({ method: 'POST', url: '/building/hiring', data, ...config })
  
  return res.data
}

/**
 * @summary Set hiring details for a building
 * {@link /building/hiring}
 */
export function usePostBuildingHiring(
  options: {
    mutation?: MutationObserverOptions<
      PostBuildingHiringMutationResponse,
      ResponseErrorConfig<PostBuildingHiring401 | PostBuildingHiring500>,
      { data?: MaybeRef<PostBuildingHiringMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBuildingHiringMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBuildingHiringMutationKey()

  return useMutation<
    PostBuildingHiringMutationResponse,
    ResponseErrorConfig<PostBuildingHiring401 | PostBuildingHiring500>,
    { data?: PostBuildingHiringMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postBuildingHiring(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}