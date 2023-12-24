import axios, { AxiosError } from 'axios'
import { useEffect, useState } from 'react'
import config from '../config'
import { IResponseEntity } from '../models'

export const useAxios = <TData>(path?: string, options: {immediate?: boolean; initialData?: any} = { immediate: true, initialData: [] }) => {
  const [response, setResponse] = useState<TData>(() => options.initialData)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [controller, setController] = useState<AbortController>()

  const computedOptions = {
    immediate: options.immediate ?? true,
    ...options
  }

  const axiosFetch = async (path: string) => {
    try {
      setLoading(true)
      const ctrl = new AbortController()
      setController(ctrl)
      const res = await axios
        .get<IResponseEntity<TData>>(config.apiBaseUrl + path, { withCredentials: true, signal: ctrl.signal })
      setResponse(res.data.data)
      return res.data.data
    } catch (err: unknown) {
      console.log(err)
      const error = err as AxiosError
      setError(error.message)
      return undefined
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    return () => controller?.abort()
  }, [controller])

  useEffect(() => {
    if (computedOptions?.immediate && path) {
      axiosFetch(path)
    }
  }, [path, options?.immediate])

  return {
    response,
    error,
    loading,
    axiosFetch
  }
}