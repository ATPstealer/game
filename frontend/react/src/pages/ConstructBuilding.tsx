import React, { useState } from 'react'
import { useParams } from 'react-router-dom'
import { Error } from '../components/Error'
import { Success } from '../components/Success'
import { constructBuilding } from '../hooks/constructBuilding'
import { useGetData } from '../hooks/getData'
import { formatDuration } from '../include/formatDuration'
import { IBuildingTypes } from '../models'

export const ConstructBuilding = () => {
  const { xParam, yParam } = useParams()
  const [x, setX] = useState(Number(xParam))
  const [y, setY] = useState(Number(yParam))
  const [typeID, setTypeID] = useState<number>(1)
  const [square, setSquare] = useState<number>(10)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const buildingTypes = useGetData<IBuildingTypes>('/building/types')
  const buildingTypeSelected = buildingTypes?.find(element => element.id === typeID)

  const submitHandler = (event: React.FormEvent) => {
    setError('')
    setMessage('')
    event.preventDefault()
    constructBuilding(x, y, typeID, square, setMessage, setError)
  }

  return (
    <>
      <div className="flex flex-col items-center justify-center mt-5 space-y-10">
        <span className="text-5xl font-bold">Building Builder</span>
        {message && <Success message={message}/> }
        {error && <Error error={error}/> }
        <form className="flex justify-center" onSubmit={submitHandler}>
          <div className="mb-4 space-y-5">
            <label htmlFor="x" className="block text-xl font-bold">Where:</label>
            <div className="flex space-x-3">
              <div className="flex flex-col">
                <label htmlFor="x" className="block text-xl font-bold">X:</label>
                <input type="number" id="x" name="x" value={x}
                  className="mt-1 border-indigo-800 border-2 rounded-md px-2"
                  onChange={(e) =>
                    setX(Number(e.target.value))}
                />
              </div>
              <div className="flex flex-col">
                <label htmlFor="y" className="block text-xl font-bold">Y:</label>
                <input type="number" id="y" name="y" value={y}
                  className="mt-1 border-indigo-800 border-2 rounded-md px-2"
                  onChange={(e) =>
                    setY(Number(e.target.value))}
                />
              </div>
            </div>
            <label htmlFor="type_id" className="block text-xl font-bold">Build:</label>
            <select id="type_id" name="type_id"
              className="mt-1 border-indigo-800 border-2 rounded-md pl-3 pr-10 text-base"
              onChange={(e) => setTypeID(Number(e.target.value))}>
              {buildingTypes?.map((buildingType, Index) => {
                return (
                  <option key={Index} value={buildingType.id}>{buildingType.title}</option>
                )
              })}
            </select>
            <div className="mt-2">{buildingTypeSelected?.description}</div>
            <div className="flex space-x-3">
              <div className="flex flex-col">
                <label htmlFor="x" className="block text-xl font-bold">Square:</label>
                <input type="number" id="square" name="square" placeholder="10"
                  className="mt-5 border-indigo-800 border-2 rounded-md pl-3"
                  onChange={(e) => setSquare(Number(e.target.value))}
                />
              </div>
            </div>
            <div className="flex space-x-3">
              <div className="flex flex-col">
                <label htmlFor="x" className="block text-xl font-bold">Cost:</label>
                {buildingTypeSelected && buildingTypeSelected.cost * square}
                <label htmlFor="x" className="block text-xl font-bold">Time:</label>
                {buildingTypeSelected && formatDuration(buildingTypeSelected.buildTime * square / 1000000000)}
              </div>
            </div>

            <button type="submit"
              className="bg-indigo-500 text-white py-2 px-4 rounded-md hover:bg-indigo-600">Construct
            </button>
          </div>
        </form>
      </div>

    </>
  )
}