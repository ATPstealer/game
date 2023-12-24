import React, { useState } from 'react'
import { useParams } from 'react-router-dom'
import { Error } from '../components/Error'
import { Success } from '../components/Success'
import { useGetData } from '../hooks/getData'
import { findResourceName } from '../include/findResourceName'
import { IGoods, IResourceTypes } from '../models'

export const Store = () => {
  const { buildingID } = useParams()
  const buildingIdNumber = Number(buildingID)
  const resourceTypes = useGetData<IResourceTypes>('/resource/types')
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const goods = useGetData<IGoods>('/store/goods/get?building_id=' + buildingID)

  return (
    <div>
      <span className="text-5xl font-bold justify-center flex mb-3">Store</span>
      {message && <Success message={message}/>}
      {error && <Error error={error}/>}

      {goods?.map((product, Index) => {
        return (
          <>
            <div className="grid grid-cols-5 gap-2 font-bold">
              <h3 className="col-span-1">Resource</h3>
              <h3 className="col-span-1">Price</h3>
              <h3 className="col-span-1">Revenue today</h3>
              <h3 className="col-span-1">Sell count today</h3>
              <h3 className="col-span-1">Status</h3>
            </div>
            <div className="grid grid-cols-5" key={Index}>
              <p className="col-span-1">{findResourceName(resourceTypes, product.resourceTypeId)}</p>
              <p className="col-span-1">{product.price}</p>
              <p className="col-span-1">{product.revenue}</p>
              <p className="col-span-1">{product.sellSum}</p>
              <p className="col-span-1">{product.status}</p>
            </div>
          </>
        )
      })}

      <div className="flex justify-center"></div>
    </div>
  )
}
