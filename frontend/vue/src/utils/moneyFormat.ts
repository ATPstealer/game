export  const moneyFormat = (amount: number| undefined): string => {
  let formatter = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    currencyDisplay: 'narrowSymbol'
  });
  return amount ? formatter.format(amount) : "Undefined";
}
