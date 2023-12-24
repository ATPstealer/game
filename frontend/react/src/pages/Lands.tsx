import React from 'react'
import { useGetData } from '../hooks/getData'
import { ILand } from '../models'

const LandsPage = () => {
  const myLand = useGetData<ILand>('/map/my')

  return (
    <div className="p-4 shadow-xl">
      <h2 className="font-bold text-2xl">Lands</h2>
      <div className="grid grid-cols-2 gap-2 font-bold">
        <h3 className="col-span-1">Cell</h3>
        <h3 className="col-span-1">Square</h3>
      </div>
      <div className="">
        {myLand?.map((land, Index) => {
          return (
            <div className="grid grid-cols-2 " key={Index}>
              <p className="col-span-1">{land.x}x{land.y}</p>
              <p className="col-span-1">{land.square} Are</p>
            </div>
          )
        })
        }
      </div>
    </div>
  )
}

export default LandsPage