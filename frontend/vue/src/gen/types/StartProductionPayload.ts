import type { TimeDuration } from './time/Duration.ts'

export interface StartProductionPayload {
  /**
   * @type integer
   */
  blueprintId: number;
  /**
   * @type string
   */
  buildingId: string;
  /**
   * @type integer
   */
  duration: TimeDuration;
}