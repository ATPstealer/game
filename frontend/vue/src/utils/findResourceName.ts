import type { ResourceType } from '@/types/Resources/index.interface'

export const findResourceName = (resourceTypes: ResourceType[] | undefined, id: number | undefined): string|undefined => {
  if (!id) {
    return ''
  }
  
  return resourceTypes?.find(resourceTypes => resourceTypes.id === id)?.name
}