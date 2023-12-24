import axios, { AxiosError } from 'axios'
import { useEffect, useState } from 'react'
import config from '../config'
import { ICellOwners, IResponseArray } from '../models'

export const useGetOwners = (
  x:number,
  y:number,
  message: string,
  error: string
) => {
  const [cellOwners, setCellOwners] = useState<ICellOwners[]>()
  const getOwners = async () => {
    try {
      const response = await axios
        .get<IResponseArray<ICellOwners>>(config.apiBaseUrl + '/map/cell_owners?x=' + x + '&y=' + y, { withCredentials: true })
      setCellOwners(response.data.data)
    } catch (e: unknown) {
      const error = e as AxiosError
      console.log(error)
    }
  }

  useEffect(() => {
    getOwners()
  }, [message, error])

  return cellOwners
}
