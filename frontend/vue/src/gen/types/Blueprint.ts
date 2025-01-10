import type { ResourceAmount } from './ResourceAmount.ts'
import type { TimeDuration } from './time/Duration.ts'

export interface Blueprint {
  /**
   * @type integer
   */
  id: number;
  /**
   * @type string
   */
  name: string;
  /**
   * @type integer
   */
  producedInId: number;
  /**
   * @type array
   */
  producedResources: ResourceAmount[];
  /**
   * @type integer
   */
  productionTime: TimeDuration;
  /**
   * @type array
   */
  usedResources: ResourceAmount[];
}