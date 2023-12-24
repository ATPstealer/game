import React, { useState } from 'react'
import { createOrder } from '../hooks/createOrder'
import { IResources } from '../models'
import { Error } from './Error'
import { Modal } from './Modal'
import { Success } from './Success'

interface SellProps {
  onClose: () => void;
  resource: IResources;
}

export const SellResource = ({ onClose, resource }: SellProps) => {
  const [amount, setAmount] = useState(0)
  const [priceForUnit, setPriceForUnit] = useState(0)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    setError('')
    setMessage('')
    await createOrder(resource.resourceTypeId, amount, priceForUnit, resource.x, resource.y,
      true, setMessage, setError)
  }

  return (
    <Modal title={'Sell ' + resource.name} onClose={onClose}>
      <Error error={error}/>
      <Success message={message}/>

      <form className="" onSubmit={handleSubmit}>
        <div className="mt-4">
          <label htmlFor="Amount" className="block text-xl font-bold">Amount:</label>
          <input name="amount" id="amount" type="number" className="border-2 w-24"
            onChange={event => {
              setAmount(parseInt(event.target.value))
            }}></input>
        </div>
        <div className="mt-4">
          <label htmlFor="Amount" className="block text-xl font-bold">Price for unit:</label>
          <input name="amount" id="amount" type="number" className="border-2 w-24"
            onChange={event => {
              setPriceForUnit(parseInt(event.target.value))
            }}></input>
        </div>

        <button type="submit"
          className="mt-4 bg-indigo-500 text-white w-1/2 py-1 px-4 rounded-md hover:bg-indigo-600"> Create Order
        </button>
      </form>

    </Modal>
  )
}
