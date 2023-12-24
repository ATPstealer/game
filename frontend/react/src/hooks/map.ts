import axios, { AxiosError } from 'axios'
import { useEffect, useState } from 'react'
import config from '../config'
import { IMap, IResponseArray } from '../models'

export const useMap = () => {
  const [map, setMap] = useState<IMap[]>()
  const [mapLoaded, setMapLoaded] = useState(false)
  const [xArray, setXArray] = useState<number[]>()
  const [yArray, setYArray] = useState<number[]>()

  const getMap = async () => {
    try {
      const response = await axios
        .get<IResponseArray<IMap>>(config.apiBaseUrl + '/map/', { withCredentials: true })
      setMap(response.data.data)
      setMapLoaded(true)
    } catch (e: unknown) {
      const error = e as AxiosError
      console.log(error)
    }
  }

  useEffect(() => {
    getMap()
  }, [])

  if (mapLoaded && !xArray && !yArray) {
    const max_x = map?.reduce((max, obj) => obj.x > max ? obj.x : max, map[0].x)
    const max_y = map?.reduce((max, obj) => obj.y > max ? obj.y : max, map[0].y)
    const min_x = map?.reduce((min, obj) => obj.x < min ? obj.x : min, map[0].x)
    const min_y = map?.reduce((min, obj) => obj.y < min ? obj.y : min, map[0].y)
    const x_array: number[] = []
    if (max_x && min_x) for (let i = min_x; i <= max_x; i++) {
      x_array.push(i)
    }
    setXArray(x_array)
    const y_array: number[] = []
    if (max_y && min_y) for (let i = max_y; i >= min_y; i--) {
      y_array.push(i)
    }
    setYArray(y_array)
  }

  return { map, mapLoaded, xArray, yArray }
}