import axios, { AxiosError } from 'axios'
import config from '../config'
import { IResponse } from '../models'

export const logout = async () => {
  try {
    await axios
      .get<IResponse>(config.apiBaseUrl + '/user/logout', { withCredentials: true })
    window.location.reload()
  } catch (e: unknown) {
    const error = e as AxiosError
    console.log(error)
  }
}