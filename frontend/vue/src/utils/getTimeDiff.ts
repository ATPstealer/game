export const getTimeDiff = (end: string | undefined): number => {
  if (!end) {
    return 0
  }
  const now = new Date()
  const workEnd = new Date(end)

  return Math.floor((workEnd.getTime() - now.getTime()) / 1000)
}