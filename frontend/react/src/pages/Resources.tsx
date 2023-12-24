import React, { useState } from 'react'
import { MoveResource } from '../components/MoveResource'
import { SellResource } from '../components/SellResource'
import { useGetData } from '../hooks/getData'
import { IResources } from '../models'

const ResourcesPage = () => {
  const myResources = useGetData<IResources>('/resource/my')
  const [moveModal, setMoveModal] = useState(false)
  const [sellModal, setSellModal] = useState(false)
  const [resource, setResource] = useState<IResources>()
  const resourceClickHandler = async (resource: IResources) => {
    setMoveModal(true)
    setResource(resource)
  }
  const sellClickHandler = async (resource: IResources) => {
    setSellModal(true)
    setResource(resource)
  }

  return (
    <>
      {moveModal && resource && <MoveResource onClose={() => {
        setMoveModal(false)
      }} resource={resource}/>}
      {sellModal && resource && <SellResource onClose={() => {
        setSellModal(false)
      }} resource={resource}/>}
      <div className="p-4 mb-4 md:mb-0 shadow-xl">
        <h2 className="font-bold text-2xl">Resources</h2>
        <div className="grid grid-cols-4 gap-2 font-bold">
          <h3 className="col-span-1">Resource</h3>
          <h3 className="col-span-1">Amount</h3>
          <h3 className="col-span-1">Cell</h3>
          <h3 className="col-span-1">Sell</h3>
        </div>
        {myResources?.map((resource, Index) => {
          return (
            resource.amount >= 1 &&
            <div className="grid grid-cols-4 " key={Index}>
              <p className="col-span-1 font-bold text-blue-500 hover:text-blue-700" onClick={() => resourceClickHandler(resource)}>{resource.name}</p>
              <p className="col-span-1">{resource.amount.toFixed(1)}</p>
              <p className="col-span-1">{resource.x}x{resource.y}</p>
              <p className="col-span-1 font-bold text-blue-500 hover:text-blue-700" onClick={() => sellClickHandler(resource)}>Sell</p>
            </div>
          )})}
        <p className="text-white"></p>
      </div>
    </>
  )
}

export default ResourcesPage