import type { JsonResult } from '@/api'

export const useUser = () => {
  const setToken = (data: JsonResult) => {
    const ttl = Number(data.data.ttl)
    const date = new Date()
    date.setTime(date.getTime() + ttl*1000)
    document.cookie = `secureToken=${data.data.token};expires=${date};domain=.${import.meta.env.VITE_DOMAIN};path=/`
    localStorage.setItem('invalid', 'false')
  }

  return {
    setToken
  }
}
