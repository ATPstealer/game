import React from 'react'
import Characteristics from '../components/Characteristics'

export const User = () => {
  return (
    <>
      <div className="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 mt-5">
        <Characteristics/>
      </div>
    </>
  )
}