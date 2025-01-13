import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'
import { computed, ref } from 'vue'

dayjs.extend(duration)

const formatTime = (time: number): string => {
  const hours = dayjs.duration(time, 'seconds').hours()
  const minutes = dayjs.duration(time, 'seconds').minutes()
  const seconds = dayjs.duration(time, 'seconds').seconds()

  return `${hours}h ${minutes}m ${seconds}s`
}

export const useTimer = () => {
  const getTime = (time: string | undefined, finish: string) => {
    if (!time) {
      return ''
    }

    const start = dayjs()
    const end = dayjs(time)

    const currentTime = ref<number>(end.diff(start, 'seconds'))
    if (currentTime.value < 0) {
      return ''
    }

    const timer = computed(() => {
      if (currentTime.value === 0) {
        return finish
      }

      return formatTime(currentTime.value)
    })

    const t = setInterval(() => {
      if (currentTime.value < 0) {
        clearInterval(t)
      }
      currentTime.value -= 1
    }, 1000)

    return timer.value
  }

  return {
    getTime
  }
}