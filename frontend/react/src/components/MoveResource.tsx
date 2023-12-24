import React, { useState } from 'react'
import { moveResource } from '../hooks/moveResource'
import { formatDuration } from '../include/formatDuration'
import { IResources } from '../models'
import { Error } from './Error'
import { Modal } from './Modal'
import { Success } from './Success'

interface MoveProps {
  onClose: () => void;
  resource: IResources;
}

export const MoveResource = ({ onClose, resource }: MoveProps) => {
  const [amount, setAmount] = useState(0)
  const [toX, setToX] = useState(Number(0))
  const [toY, setToY] = useState(Number(0))
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const distance = ((resource.x-toX)**2 + (resource.y-toY)**2)**(0.5)
  const price = (resource.weight + resource.volume) * distance * amount / 1000

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    setError('')
    setMessage('')
    await moveResource(resource.resourceTypeId, amount, resource.x, resource.y, toX, toY, setMessage, setError)
  }

  return (
    <Modal title="Move Resources" onClose={onClose}>
      <Error error={error}/>
      <Success message={message}/>
      <div className="text-xl">
        <span className="font-bold">{resource.name}</span> move from <span
          className="font-bold">{resource.x}:{resource.y}</span> to
      </div>
      <form className="" onSubmit={handleSubmit}>
        <div className="flex space-x-3">
          <div className="flex flex-col">
            <label htmlFor="x" className="block text-xl font-bold">X:</label>
            <input type="number" id="x" name="x" value={toX}
              className="border-indigo-800 border-2 rounded-md px-2 w-16"
              onChange={(e) =>
                setToX(Number(e.target.value))}
            />
          </div>
          <div className="flex flex-col">
            <label htmlFor="y" className="block text-xl font-bold">Y:</label>
            <input type="number" id="y" name="y" value={toY}
              className="border-indigo-800 border-2 rounded-md px-2 w-16"
              onChange={(e) =>
                setToY(Number(e.target.value))}
            />
          </div>
        </div>
        <div className="mt-4">
          <label htmlFor="Amount" className="block text-xl font-bold">Amount:</label>
          <input name="amount" id="amount" type="number" className="border-2 w-24"
            onChange={event => {
              setAmount(parseInt(event.target.value))
            }}></input>
        </div>
        <div className="mt-4 font-bold">Estimate price: {price} </div>
        <div className="font-bold">Estimate time: {formatDuration(distance*600)} </div>
        <button type="submit"
          className="mt-4 bg-indigo-500 text-white w-1/2 py-1 px-4 rounded-md hover:bg-indigo-600"> Move
        </button>
      </form>
    </Modal>
  )
}
