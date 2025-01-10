export interface Bank {
  /**
   * @type number
   */
  borrowedFromState: number;
  /**
   * @type number
   */
  borrowedLimit: number;
  /**
   * @type number
   */
  loansAmount: number;
  /**
   * @type number
   */
  loansAmountNewUsers: number;
  /**
   * @type number
   */
  loansLimit: number;
}