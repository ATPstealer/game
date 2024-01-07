export  const moneyFormat = (amount: number| undefined): string => {
  const formatter = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    currencyDisplay: 'narrowSymbol'
  })

  return amount ? formatter.format(amount) : 'Undefined'
}
