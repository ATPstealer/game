import type { TimeDuration } from './time/Duration.ts'

export interface BuildingType {
  /**
   * @type integer
   */
  buildTime: TimeDuration;
  /**
   * @type string
   */
  buildingGroup: string;
  /**
   * @type string
   */
  buildingSubGroup: string;
  /**
   * @type number
   */
  capacity: number;
  /**
   * @type number
   */
  cost: number;
  /**
   * @type string
   */
  description: string;
  /**
   * @type integer
   */
  id: number;
  /**
   * @type string
   */
  requirements: string;
  /**
   * @type string
   */
  title: string;
  /**
   * @type integer
   */
  workers: number;
}