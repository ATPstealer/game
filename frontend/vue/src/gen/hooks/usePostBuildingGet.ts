import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { PostBuildingGetMutationRequest, PostBuildingGetMutationResponse, PostBuildingGet500 } from '../types/PostBuildingGet.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig, ResponseErrorConfig } from '@/api/customClientAxios'

export const postBuildingGetMutationKey = () => [{ url: '/building/get' }] as const

export type PostBuildingGetMutationKey = ReturnType<typeof postBuildingGetMutationKey>

/**
 * @summary Fetch the list of buildings
 * {@link /building/get}
 */
async function postBuildingGet(data?: PostBuildingGetMutationRequest, config: Partial<RequestConfig<PostBuildingGetMutationRequest>> = {}) {
  const res = await client<PostBuildingGetMutationResponse, ResponseErrorConfig<PostBuildingGet500>, PostBuildingGetMutationRequest>({
    method: 'POST',
    url: '/building/get',
    data,
    ...config
  })
  
  return res.data
}

/**
 * @summary Fetch the list of buildings
 * {@link /building/get}
 */
export function usePostBuildingGet(
  options: {
    mutation?: MutationObserverOptions<
      PostBuildingGetMutationResponse,
      ResponseErrorConfig<PostBuildingGet500>,
      { data?: MaybeRef<PostBuildingGetMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBuildingGetMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBuildingGetMutationKey()

  return useMutation<PostBuildingGetMutationResponse, ResponseErrorConfig<PostBuildingGet500>, { data?: PostBuildingGetMutationRequest }>({
    mutationFn: async ({ data }) => {
      return postBuildingGet(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}