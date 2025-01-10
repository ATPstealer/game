import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBuildingStopWorkMutationRequest,
  PostBuildingStopWorkMutationResponse,
  PostBuildingStopWork401,
  PostBuildingStopWork500
} from '../types/PostBuildingStopWork.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const postBuildingStopWorkMutationKey = () => [{ url: '/building/stop_work' }] as const

export type PostBuildingStopWorkMutationKey = ReturnType<typeof postBuildingStopWorkMutationKey>

/**
 * @summary Stops any work in building. Later he should stop only the works available for stopping.
 * {@link /building/stop_work}
 */
async function postBuildingStopWork(data: PostBuildingStopWorkMutationRequest, config: Partial<RequestConfig<PostBuildingStopWorkMutationRequest>> = {}) {
  const res = await client<PostBuildingStopWorkMutationResponse, PostBuildingStopWork401 | PostBuildingStopWork500, PostBuildingStopWorkMutationRequest>({
    method: 'POST',
    url: '/building/stop_work',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @summary Stops any work in building. Later he should stop only the works available for stopping.
 * {@link /building/stop_work}
 */
export function usePostBuildingStopWork(
  options: {
    mutation?: MutationObserverOptions<
      PostBuildingStopWorkMutationResponse,
      PostBuildingStopWork401 | PostBuildingStopWork500,
      { data: MaybeRef<PostBuildingStopWorkMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBuildingStopWorkMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBuildingStopWorkMutationKey()

  return useMutation<PostBuildingStopWorkMutationResponse, PostBuildingStopWork401 | PostBuildingStopWork500, { data: PostBuildingStopWorkMutationRequest }>({
    mutationFn: async ({ data }) => {
      return postBuildingStopWork(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}