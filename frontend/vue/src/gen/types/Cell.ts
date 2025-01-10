export interface Cell {
  /**
   * @type string
   */
  _id: string;
  /**
   * @type number
   */
  SpendRate: number;
  /**
   * @type number
   */
  averageSalary: number;
  /**
   * @type string
   */
  cellName: string;
  /**
   * @type number | undefined
   */
  civilSavings?: number;
  /**
   * @type number
   */
  crime: number;
  /**
   * @type number
   */
  education: number;
  /**
   * @type number
   */
  medicine: number;
  /**
   * @type number
   */
  pollution: number;
  /**
   * @type number
   */
  population: number;
  /**
   * @type integer
   */
  square: number;
  /**
   * @type string
   */
  surfaceImagePath: string;
  /**
   * @type integer
   */
  x: number;
  /**
   * @type integer
   */
  y: number;
}