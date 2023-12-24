import { Button } from 'primereact/button'
import { Dropdown } from 'primereact/dropdown'
import React, { useState } from 'react'
import { useParams } from 'react-router-dom'
import { timeValues } from '../components/Buildings/constants'
import ResourceCard from '../components/Buildings/ResourceCard'
import { Error } from '../components/Error'
import { Success } from '../components/Success'
import { destroyBuilding } from '../hooks/destroyBuilding'
import { useGetData } from '../hooks/getData'
import { startWork } from '../hooks/startWork'
import { IBlueprint, IBuildings, IResourceTypes } from '../models'

export const Building = () => {
  const { buildingID } = useParams()
  const buildingIdNumber = Number(buildingID)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const resourceTypes = useGetData<IResourceTypes>('/resource/types')
  const myBuildings = useGetData<IBuildings>('/building/my')
  const blueprintsAll = useGetData<IBlueprint>('/building/blueprints')
  const building = myBuildings?.find(myBuildings => myBuildings.id === buildingIdNumber)
  const blueprintsForBuilding = blueprintsAll?.filter(blueprintsAll => blueprintsAll.producedInId === building?.typeId)
  const [selectedBlueprint, setSelectedBlueprint] = useState<number | null>(null)
  const [selectedDuration, setSelectedDuration] = useState(timeValues[1].value)
  const submitHandler = (event: React.FormEvent) => {
    setError('')
    setMessage('')
    event.preventDefault()
    if (selectedBlueprint) {
      startWork(buildingIdNumber, selectedBlueprint, selectedDuration, setMessage, setError)
    } else {
      setError('Choose resource for production')
    }
  }

  const pickBlueprint = (blueprint: number) => {
    setSelectedBlueprint(blueprint)
  }

  return (
    <div>
      <span className="text-5xl font-bold justify-center flex mb-3">Production</span>
      {message && <Success message={message}/>}
      {error && <Error error={error}/>}

      <div className="flex justify-center">
        <div  className="space-y-4 w-full max-w-md bg-white p-6 rounded-xl shadow-md">
          <div className="font-bold">Start work:</div>
          <Dropdown
            value={selectedDuration}
            onChange={event => setSelectedDuration(event.value)}
            options={timeValues}
            optionLabel="label"
            className="w-full"
          />
          <Button
            disabled={!selectedBlueprint}
            onClick={submitHandler}
            label="Start work"
            className="w-full"
          />
        </div>
        <div className="space-y-4 w-full max-w-md bg-white p-6 rounded-xl shadow-md">
          <div className="font-bold">Building:</div>
          Type: {building?.title} <br/>
          Status: {building?.status} <br/>
          Coordinates: {building?.x}:{building?.y} <br/>
          Level x Square: {building?.level}x{building?.square} <br/>
          {building?.blueprintName && <>Producing now: {building.blueprintName} <br/> </>}
          {building?.status !== 'Ready' && <>Finish work : {building?.workEnd} <br/> </>}
          <div className="text-red-500 hover:text-red-700"
            onClick={() => destroyBuilding(building?.id, setMessage, setError)}>One click = destroy building
          </div>
        </div>
      </div>
      <div className="mt-3 flex font-bold justify-center">Choose resource for creating:</div>
      <div className="mt-5 mx-5 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {
          building &&
          blueprintsForBuilding &&
          blueprintsAll &&
          blueprintsForBuilding.map(blueprint =>
            <ResourceCard key={blueprint.id} {...{ blueprint, resourceTypes, pickBlueprint, selectedBlueprint }}/>
          )
        }
      </div>
    </div>
  )
}