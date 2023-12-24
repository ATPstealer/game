import React from 'react'
import { useGetData } from '../hooks/getData'
import { IStorage } from '../models'

const StoragesPage = () => {
  const myStorages = useGetData<IStorage>('/storage/my')

  return (
    <div className="p-4 mb-4 md:mb-0 shadow-xl">
      <h2 className="font-bold text-2xl">Storages</h2>
      <div className="grid grid-cols-3 gap-2 font-bold">
        <h3 className="col-span-1">Cell</h3>
        <h3 className="col-span-1">Volumes</h3>
      </div>
      {myStorages?.map((storage, Index) => {
        return (
          <div className="grid grid-cols-3 " key={Index}>
            <p className="col-span-1">{storage.x}:{storage.y}</p>
            <p className={storage.volumeOccupied > storage.volumeMax ? 'font-bold text-red-600' : ''}>
              {storage.volumeOccupied}/{storage.volumeMax}
            </p>
          </div>
        )})}
      <p className="text-white"></p>
    </div>
  )
}

export default StoragesPage