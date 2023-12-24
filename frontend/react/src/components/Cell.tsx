import axios, { AxiosError } from 'axios'
import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import config from '../config'
import { useGetOwners } from '../hooks/getOwners'
import { ICellOwners, IMap, IResponseArray } from '../models'
import { Error } from './Error'
import { Modal } from './Modal'
import { Success } from './Success'

interface CellProps {
  onClose: () => void;
  cell: IMap;
  square: number;
}

export const Cell = ({ onClose, cell, square }: CellProps) => {
  const [buySquare, setBuySquare] = useState<number>()
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const cellOwners = useGetOwners(cell.x, cell.y, message, error)
  let landOccupied: number = 0
  if (cellOwners) {
    for (const owner of cellOwners) {
      landOccupied += owner.square
    }
  }

  let freeSquare = square
  if (cellOwners) {
    for (const owner of cellOwners) {
      freeSquare -= owner.square
    }
  }

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    setError('')
    setMessage('')

    if (typeof buySquare === 'undefined') return
    try {
      const response = await axios
        .get<IResponseArray<ICellOwners>>(config.apiBaseUrl +
          '/map/buy_land?x=' + cell.x + '&y=' + cell.y + '&square=' + buySquare, { withCredentials: true })
      console.log(response.data)
      if (response.data.status === 'success') {
        setMessage(response.data.text)
      } else {
        setError(response.data.text)
      }
    } catch (e: unknown) {
      const error = e as AxiosError
      console.log(error)
    }
  }

  return (
    <>
      <Modal title={`Cell ${cell.x} x ${cell.y}`} onClose={onClose}>
        <Error error={error}/>
        <Success message={message}/>
        <p>Each free slot of land costs $10 more than the previous one.
          Thus, the first one costs $10, the hundredth one costs $1000</p><br/>
        <div className="font-bold">Free space: {freeSquare} Ares</div>
        <div className="font-bold">Occupied land: {landOccupied} Ares</div>
        <form className="mb-4" onSubmit={handleSubmit}>
          <p> Buy square: <input name="bay_square" id="bay_square" type="number" className="border-2 w-20"
            onChange={event => {
              setBuySquare(parseInt(event.target.value))
            }}></input>  </p>
          {buySquare && <p> Estimate price: {10 * (landOccupied*2 + 1 + buySquare ) * buySquare / 2} $
            <p>
              <button type="submit"
                className="bg-indigo-500 text-white w-1/2 py-1 px-4 rounded-md hover:bg-indigo-600"> Buy </button>
            </p>
          </p>}
        </form>
        <Link to={`/construct_building/${cell.x}/${cell.y}`} className="font-bold text-blue-500 hover:text-blue-700"> Construct building</Link>

        <div className="font-bold mt-4"> Pollution: {cell.pollution}</div>
        <div className="font-bold"> Population: {cell.population}</div>
        <div className="font-bold"> Education: {cell.education}</div>
        <div className="font-bold"> Crime: {cell.crime}</div>
        <div className="font-bold"> Medicine: {cell.medicine}</div>
        <div className="font-bold"> ElementarySchool: {cell.elementarySchool}</div>
        <div className="font-bold"> HigherSchool: {cell.higherSchool}</div>
        <div className="font-bold"> Landlords:</div>
        {
          cellOwners?.map((owner, index) => {
            return (<p key={index}>{owner.nickName}: {owner.square} Are</p>)
          })
        }
      </Modal>
    </>
  )
}
