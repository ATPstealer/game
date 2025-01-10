import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  DeleteBuildingDestroyMutationResponse,
  DeleteBuildingDestroyQueryParams,
  DeleteBuildingDestroy401,
  DeleteBuildingDestroy500
} from '../types/DeleteBuildingDestroy.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const deleteBuildingDestroyMutationKey = () => [{ url: '/building/destroy' }] as const

export type DeleteBuildingDestroyMutationKey = ReturnType<typeof deleteBuildingDestroyMutationKey>

/**
 * @summary Destroy an existing building
 * {@link /building/destroy}
 */
async function deleteBuildingDestroy(params: DeleteBuildingDestroyQueryParams, config: Partial<RequestConfig> = {}) {
  const res = await client<DeleteBuildingDestroyMutationResponse, DeleteBuildingDestroy401 | DeleteBuildingDestroy500, unknown>({
    method: 'DELETE',
    url: '/building/destroy',
    baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2',
    params,
    ...config
  })
  
  return res.data
}

/**
 * @summary Destroy an existing building
 * {@link /building/destroy}
 */
export function useDeleteBuildingDestroy(
  options: {
    mutation?: MutationObserverOptions<
      DeleteBuildingDestroyMutationResponse,
      DeleteBuildingDestroy401 | DeleteBuildingDestroy500,
      { params: MaybeRef<DeleteBuildingDestroyQueryParams> }
    >;
    client?: Partial<RequestConfig>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? deleteBuildingDestroyMutationKey()

  return useMutation<DeleteBuildingDestroyMutationResponse, DeleteBuildingDestroy401 | DeleteBuildingDestroy500, { params: DeleteBuildingDestroyQueryParams }>({
    mutationFn: async ({ params }) => {
      return deleteBuildingDestroy(params, config)
    },
    mutationKey,
    ...mutationOptions
  })
}