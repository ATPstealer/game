import React from 'react'

interface SuccessProps {
  message: string;
}

export const Success = ({ message }: SuccessProps) => {
  return (
    <p className="text-center text-green-600 font-bold mb-2">{ message }</p>
  )
}
