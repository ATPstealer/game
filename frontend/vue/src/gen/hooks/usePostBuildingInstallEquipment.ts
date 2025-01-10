import type { MutationObserverOptions } from '@tanstack/vue-query'
import { useMutation } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type {
  PostBuildingInstallEquipmentMutationRequest,
  PostBuildingInstallEquipmentMutationResponse,
  PostBuildingInstallEquipment401,
  PostBuildingInstallEquipment500
} from '../types/PostBuildingInstallEquipment.ts'
import client from '@/api/customClientAxios'
import type { RequestConfig } from '@/api/customClientAxios'

export const postBuildingInstallEquipmentMutationKey = () => [{ url: '/building/install_equipment' }] as const

export type PostBuildingInstallEquipmentMutationKey = ReturnType<typeof postBuildingInstallEquipmentMutationKey>

/**
 * @summary Install equipment in a building
 * {@link /building/install_equipment}
 */
async function postBuildingInstallEquipment(
  data: PostBuildingInstallEquipmentMutationRequest,
  config: Partial<RequestConfig<PostBuildingInstallEquipmentMutationRequest>> = {}
) {
  const res = await client<
    PostBuildingInstallEquipmentMutationResponse,
    PostBuildingInstallEquipment401 | PostBuildingInstallEquipment500,
    PostBuildingInstallEquipmentMutationRequest
  >({ method: 'POST', url: '/building/install_equipment', baseURL: 'http://staging.game.k8s.atpstealer.com/api/v2', data, ...config })
  
  return res.data
}

/**
 * @summary Install equipment in a building
 * {@link /building/install_equipment}
 */
export function usePostBuildingInstallEquipment(
  options: {
    mutation?: MutationObserverOptions<
      PostBuildingInstallEquipmentMutationResponse,
      PostBuildingInstallEquipment401 | PostBuildingInstallEquipment500,
      { data: MaybeRef<PostBuildingInstallEquipmentMutationRequest> }
    >;
    client?: Partial<RequestConfig<PostBuildingInstallEquipmentMutationRequest>>;
  } = {}
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? postBuildingInstallEquipmentMutationKey()

  return useMutation<
    PostBuildingInstallEquipmentMutationResponse,
    PostBuildingInstallEquipment401 | PostBuildingInstallEquipment500,
    { data: PostBuildingInstallEquipmentMutationRequest }
  >({
    mutationFn: async ({ data }) => {
      return postBuildingInstallEquipment(data, config)
    },
    mutationKey,
    ...mutationOptions
  })
}