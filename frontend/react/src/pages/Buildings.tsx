import React from 'react'
import { useGetData } from '../hooks/getData'
import { formatDuration } from '../include/formatDuration'
import { IBuildings } from '../models'

const BuildingsPage = () => {
  const myBuildings = useGetData<IBuildings>('/building/my')

  return (
    <div className="p-4 mb-4 md:mb-0 shadow-xl">
      <h2 className="font-bold text-2xl">Buildings
        <a className="text-blue-500 hover:text-blue-700"
          href="/construct_building/0/0"> Construct</a>
      </h2>
      <div className="text-black">
        <div className="grid grid-cols-4 gap-2 font-bold">
          <h3 className="col-span-1">Building</h3>
          <h3 className="col-span-1">Cell</h3>
          <h3 className="col-span-1">Status</h3>
          <h3 className="col-span-1">Finish</h3>
        </div>
        {myBuildings?.map((building, Index) => {
          let diffS: number | null
          building.workEnd = new Date(building.workEnd)
          const now = new Date()
          diffS = Math.floor((building.workEnd.getTime() - now.getTime()) / 1000)
          if (diffS < 0) diffS = null
          return (
            <div className="grid grid-cols-4 " key={Index}>
              <p className="col-span-1 font-bold text-blue-500 hover:text-blue-700">
                <a href={`/building/${building.id}`}>{building.title} {building.level}x{building.square}</a></p>
              <p className="col-span-1">{building.x}x{building.y}</p>
              <p className="col-span-1">{building.status}</p>
              <p className="col-span-1 ml-2">{diffS && formatDuration(diffS)}</p>
            </div>
          )
        })
        }
      </div>
    </div>
  )
}

export default BuildingsPage