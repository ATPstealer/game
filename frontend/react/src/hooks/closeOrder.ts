import axios, { AxiosError } from 'axios'
import config from '../config'
import { IResponse, Message } from '../models'

// v1/market/order/close?order_id=15
export const closeOrder = async (
  orderID: number,
  setMessage: Message,
  setError: Message
) => {
  try {
    const response = await axios
      .get<IResponse>(config.apiBaseUrl +
        '/market/order/close?order_id=' + orderID, { withCredentials: true })
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