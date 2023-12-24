import { IResourceTypes } from '../models'

export const findResourceName = (resourceTypes: IResourceTypes[]|undefined, id: number): string|undefined => {
  return resourceTypes?.find(resourceTypes => resourceTypes.id === id)?.name
}