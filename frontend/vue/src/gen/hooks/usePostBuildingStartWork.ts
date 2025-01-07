import client from '@kubb/plugin-client/clients/fetch'
import type { RequestConfig } from '@kubb/plugin-client/clients/fetch'
import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBuildingStartWorkMutationRequest,
  PostBuildingStartWorkMutationResponse,
  PostBuildingStartWork401,
  PostBuildingStartWork500
} from '../types/PostBuildingStartWork.ts'

export const postBuildingStartWorkMutationKey = () => [{ url: '/building/start_work' }] as const

export type PostBuildingStartWorkMutationKey = ReturnType<typeof postBuildingStartWorkMutationKey>

/**
 * @summary Start work in the building
 * {@link /building/start_work}
 */
async function postBuildingStartWork(data: PostBuildingStartWorkMutationRequest, config: Partial<RequestConfig<PostBuildingStartWorkMutationRequest>> = {}) {
  const res = await client<PostBuildingStartWorkMutationResponse, PostBuildingStartWork401 | PostBuildingStartWork500, PostBuildingStartWorkMutationRequest>({
    method: 'POST',
    url: '/building/start_work',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @summary Start work in the building
 * {@link /building/start_work}
 */
export function usePostBuildingStartWork(
  options: {
    mutation?: MutationObserverOptions<
      PostBuildingStartWorkMutationResponse,
      PostBuildingStartWork401 | PostBuildingStartWork500,
      { data: MaybeRef<PostBuildingStartWorkMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBuildingStartWorkMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBuildingStartWorkMutationKey()

  return useMutation<
    PostBuildingStartWorkMutationResponse,
    PostBuildingStartWork401 | PostBuildingStartWork500,
    { data: PostBuildingStartWorkMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postBuildingStartWork(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}