import React from 'react'

interface ModalProps {
  children: React.ReactNode;
  title: string;
  onClose: () => void;
}

export const Modal = ({ children, title, onClose }: ModalProps) => {
  return (
    <>
      <div className="fixed bg-black/50 top-0 right-0 left-0 bottom-0 z-[1]" onClick={onClose}/>
      <div className="w-[400px] sm:w-full p-5 rounded bg-white fixed lg:top-28 left-1/2 -translate-x-1/2 shadow-md  max-w-md z-[2]">
        <h2 className="text-3xl font-semibold mb-6">{title}</h2>
        {children}
      </div>
    </>
  )
}