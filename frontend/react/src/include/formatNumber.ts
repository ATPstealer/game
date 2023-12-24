export const formatNumber = (n: number): string => {
  const nf = new Intl.NumberFormat('fr-CH')
  return nf.format(n)
}