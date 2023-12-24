import axios, { AxiosError } from 'axios'
import config from '../config'
import { IResponse, Message } from '../models'

// v1/market/order/create?resource_type_id=3&x=1&y=0&amount=100&price_for_unit=100&sell=true
export const createOrder = async (
  resourceTypeID: number,
  amount: number,
  priceForUnit: number,
  x: number,
  y: number,
  sell: boolean,
  setMessage: Message,
  setError: Message
) => {
  try {
    const response = await axios
      .get<IResponse>(config.apiBaseUrl +
        '/market/order/create?resource_type_id=' + resourceTypeID + '&amount=' + amount + '&x=' + x +
        '&y=' + y + '&price_for_unit=' + priceForUnit + '&sell=' + sell.toString(), { withCredentials: true })
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
