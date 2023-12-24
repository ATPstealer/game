import axios, { AxiosError } from 'axios'
import { useEffect, useState } from 'react'
import config from '../config'
import { IResponseEntity, IResponseArray } from '../models'

export const useGetData = <TData>(path: string) => {
  const [data, setData] = useState<TData[]>()

  const getData = async () => {
    try {
      const response = await axios
        .get<IResponseArray<TData>>(config.apiBaseUrl + path, { withCredentials: true })
      setData(response.data.data)
    } catch (e: unknown) {
      const error = e as AxiosError
      console.log(error)
    }
  }

  useEffect(() => {
    getData()
  }, [path])

  return data || []
}

export const useGetEntity = <TData>(path: string) => {
  const [data, setData] = useState<TData>()

  const getData = async () => {
    try {
      const response = await axios
        .get<IResponseEntity<TData>>(config.apiBaseUrl + path, { withCredentials: true })
      setData(response.data.data)
    } catch (e: unknown) {
      const error = e as AxiosError
      console.log(error)
    }
  }

  useEffect(() => {
    getData()
  }, [])

  return data || {} as TData
}