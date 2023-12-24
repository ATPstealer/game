import axios, { AxiosError } from 'axios'
import config from '../config'
import { IResponse, Message } from '../models'

export const startWork = async (
  buildingID:number,
  blueprintID:number,
  duration:number,
  setMessage: Message,
  setError: Message
) => {
  try {
    const response = await axios
      .get<IResponse>(config.apiBaseUrl +
        '/building/start_work?building_id=' + buildingID + '&blueprint_id=' + blueprintID + '&duration=' + duration, { withCredentials: true })
    if (response.data.status === 'success') {
      setMessage(response.data.text)
    } else {
      setError(response.data.text)
    }
  } catch (e: unknown) {
    const error = e as AxiosError
    setError(error.message)
    console.log(error)
  }
}
