import type { ResourceType } from '@/types/Resources/index.interface'

export const findResourceName = (resourceTypes: ResourceType[]|undefined, id: number): string|undefined => {
  return resourceTypes?.find(resourceTypes => resourceTypes.id === id)?.name
}