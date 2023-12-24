import React from 'react'
import { useGetData } from '../hooks/getData'
import { formatDuration } from '../include/formatDuration'
import { ILogistics } from '../models'

const LogisticsPage = () => {
  const myLogistics = useGetData<ILogistics>('/resource/my_logistics')

  return (
    <div className="p-4 shadow-xl">
      <h2 className="font-bold text-2xl">Logistics</h2>
      <div className="grid grid-cols-5 font-bold">
        <h3 className="col-span-1">Resource</h3>
        <h3 className="col-span-1">From </h3>
        <h3 className="col-span-1">To </h3>
        <h3 className="col-span-1">Amount</h3>
        <h3 className="col-span-1">Finish</h3>
      </div>
      <div className="">
        {myLogistics?.map((logistic, Index) => {
          let diffS: number | null
          const now = new Date()
          logistic.workEnd = new Date(logistic.workEnd)
          diffS = Math.floor((logistic.workEnd.getTime() - now.getTime()) / 1000)
          if (diffS < 0) diffS = null
          return (
            <div className="grid grid-cols-5 " key={Index}>
              <p className="col-span-1">{logistic.resourceName}</p>
              <p className="col-span-1">{logistic.fromX}x{logistic.fromY}</p>
              <p className="col-span-1">{logistic.toX}x{logistic.toY}</p>
              <p className="col-span-1">{logistic.amount}</p>
              <p className="col-span-1">{diffS && formatDuration(diffS)}</p>
            </div>
          )
        })
        }
      </div>
    </div>
  )
}

export default LogisticsPage