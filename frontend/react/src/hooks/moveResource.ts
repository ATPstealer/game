import axios, { AxiosError } from 'axios'
import config from '../config'
import { IResponse, Message } from '../models'

export const moveResource = async (
  resourceTypeID: number,
  amount: number,
  fromX: number,
  fromY: number,
  toX: number,
  toY: number,
  setMessage: Message,
  setError: Message
) => {
  try {
    const response = await axios
      .get<IResponse>(config.apiBaseUrl +
        '/resource/move?resource_type_id=' + resourceTypeID + '&amount=' + amount +
        '&from_x=' + fromX+ '&from_y=' + fromY + '&to_x=' + toX + '&to_y=' + toY, { withCredentials: true })
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
