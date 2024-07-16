import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'

dayjs.extend(duration)

export const formatDuration = (time: number): string => {
  const hours = dayjs.duration(time, 'seconds').hours()
  const minutes = dayjs.duration(time, 'seconds').minutes()
  const seconds = dayjs.duration(time, 'seconds').seconds()

  return `${hours}h ${minutes}m ${seconds}s`
}