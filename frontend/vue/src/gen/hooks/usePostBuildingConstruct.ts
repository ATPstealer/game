import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBuildingConstructMutationRequest,
  PostBuildingConstructMutationResponse,
  PostBuildingConstruct401,
  PostBuildingConstruct500
} from '../types/PostBuildingConstruct.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const postBuildingConstructMutationKey = () => [{ url: '/building/construct' }] as const

export type PostBuildingConstructMutationKey = ReturnType<typeof postBuildingConstructMutationKey>

/**
 * @summary Construct a new building
 * {@link /building/construct}
 */
async function postBuildingConstruct(data: PostBuildingConstructMutationRequest, config: Partial<RequestConfig<PostBuildingConstructMutationRequest>> = {}) {
  const res = await client<
    PostBuildingConstructMutationResponse,
    ResponseErrorConfig<PostBuildingConstruct401 | PostBuildingConstruct500>,
    PostBuildingConstructMutationRequest
  >({ method: 'POST', url: '/building/construct', data, ...config })
  
  return res.data
}

/**
 * @summary Construct a new building
 * {@link /building/construct}
 */
export function usePostBuildingConstruct(
  options: {
    mutation?: MutationObserverOptions<
      PostBuildingConstructMutationResponse,
      ResponseErrorConfig<PostBuildingConstruct401 | PostBuildingConstruct500>,
      { data: MaybeRef<PostBuildingConstructMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBuildingConstructMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBuildingConstructMutationKey()

  return useMutation<
    PostBuildingConstructMutationResponse,
    ResponseErrorConfig<PostBuildingConstruct401 | PostBuildingConstruct500>,
    { data: PostBuildingConstructMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postBuildingConstruct(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}