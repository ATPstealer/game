import React, { useState } from 'react'
import { Cell } from '../components/Cell'
import config from '../config'
import { useMap } from '../hooks/map'
import { IMap } from '../models'

export const Map = () => {
  const { map, mapLoaded, xArray, yArray } = useMap()
  const [cellModal, setCellModal] = useState(false)
  const [cell, setCell] = useState<IMap>()
  const [square, setSquare] = useState<number>(0)

  const CellClickHandler = async (oneCell: IMap) => {
    setCellModal(true)
    setCell(oneCell)
    setSquare(oneCell.square)
  }

  return (
    <>
      {cellModal && cell && <Cell cell={cell} square={square}
        onClose={() => {
          setCellModal(false)
        }}
      />}
      {mapLoaded && xArray && yArray &&
        yArray.map((row, rowIndex) => {
          const itemsInRow = map?.filter(item => item.y === row)
          return (
            <div className="flex justify-center" key={rowIndex}>
              {xArray.map((column, columnIndex) => {
                const oneCell = itemsInRow?.filter(item => item.x === column)[0]
                if (oneCell) {
                  return (
                    <div className="relative" key={columnIndex}>
                      <img className="border border-blue-950" alt="Surface"
                        src={config.minioUrl + oneCell?.surfaceImagePath}
                        onClick={() => {
                          CellClickHandler(oneCell)
                        }}/>
                      <p className="absolute top-0 left-0 text-white">{oneCell?.cellName}</p>
                    </div>
                  )
                }
              })}
            </div>
          )
        })}
    </>
  )
}
