import type { Characteristics } from './Characteristics.ts'

export interface User {
  /**
   * @type string
   */
  _id: string;
  /**
   * @type object | undefined
   */
  characteristics?: Characteristics;
  /**
   * @type string | undefined
   */
  created?: string;
  /**
   * @type number | undefined
   */
  creditRating?: number;
  /**
   * @type string
   */
  email: string;
  /**
   * @type number | undefined
   */
  money?: number;
  /**
   * @type string
   */
  nickName: string;
  /**
   * @type string
   */
  password: string;
}