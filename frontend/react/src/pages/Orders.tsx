import React from 'react'
import { closeOrder } from '../hooks/closeOrder'
import { useGetData } from '../hooks/getData'
import { IOrder } from '../models'

const OrdersPage = () => {
  const myOrders = useGetData<IOrder>('/market/order/my')

  return (
    <div className="p-4 mb-4 md:mb-0 shadow-xl">
      <h2 className="font-bold text-2xl">Orders</h2>
      <div className="grid grid-cols-5 gap-2 font-bold">
        <h3 className="col-span-1">Cell</h3>
        <h3 className="col-span-1">Resource</h3>
        <h3 className="col-span-1">Amount</h3>
        <h3 className="col-span-1">Price</h3>
        <h3 className="col-span-1">Type (Close)</h3>
      </div>
      {myOrders?.map((order, Index) => {
        return (
          <div className="grid grid-cols-5 " key={Index}>
            <p className="col-span-1">{order.x}:{order.y}</p>
            <p className="col-span-1">{order.resourceName}</p>
            <p className="col-span-1">{order.amount}</p>
            <p className="col-span-1">{order.priceForUnit}$</p>
            <p className="col-span-1 text-red-600 font-bold">
              <span onClick={() => closeOrder(order.id, () => {}, () => {})}>
                {order.sell ? 'sell' : 'buy'}</span></p>
          </div>
        )})}
      <p className="text-white"></p>
    </div>
  )
}

export default OrdersPage