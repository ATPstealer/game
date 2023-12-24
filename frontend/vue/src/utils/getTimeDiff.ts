export const getTimeDiff = (end: string): number => {
  const now = new Date()
  const workEnd = new Date(end)

  return Math.floor((workEnd.getTime() - now.getTime()) / 1000)
}