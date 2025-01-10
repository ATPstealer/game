export interface JsonResult {
  /**
   * @type integer
   */
  code: number;
  data?: unknown;
  /**
   * @type string | undefined
   */
  text?: string;
  values?: any;
}