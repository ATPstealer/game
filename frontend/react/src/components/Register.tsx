import axios, { AxiosError } from 'axios'
import sha256 from 'crypto-js/sha256'
import React, { useState } from 'react'
import config from '../config'
import { IUser, IResponseArray } from '../models'
import { Error } from './Error'
import { Modal } from './Modal'
import { Success } from './Success'

interface RegisterProps {
  onClose: () => void;
  switchToLogin: () => void;
}

const userData: IUser = {
  nickName: '',
  email: '',
  password: ''
}

export const Register = ({ onClose, switchToLogin }: RegisterProps) => {
  const [nickName, setNickName] = useState('') // TODO: как сделать чтобы призакрытии модалки данные жили
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  const submitHandler = async (event: React.FormEvent) => {
    event.preventDefault()
    setError('')
    setMessage('')

    // Check Fields
    if (nickName.length < 3) {
      setError('Fill in the field nickname')
      return
    }
    if (email.length < 5) {
      setError('Fill in the field email')
      return
    }
    if (!isStrongPassword(password)) {
      setError('Password must include lowercase and uppercase later, digits and be more than 8 symbols')
      return
    }
    if (password !== confirmPassword) {
      setError('Passwords doesn\'t match')
      return
    }

    // Try to register user
    userData.nickName = nickName
    userData.email = email
    userData.password = sha256(password).toString()
    try {
      const response = await axios.post<IResponseArray<IUser>>(config.apiBaseUrl + '/user/create', userData)
      if (response.data.status === 'success') {
        setMessage('User successfully registered. Please, login')
      } else {
        setError('User didn\'t registered: ' + response.data.text)
      }
    } catch (e: unknown) {
      const error = e as AxiosError
      setError(error.message)
    }
  }

  return (
    <Modal title="Sing up" onClose={onClose}>
      <Error error={error} />
      <Success message={message} />
      <form onSubmit={submitHandler}>
        <div className="mb-4">
          <label htmlFor="nickname" className="block text-gray-700 text-sm font-medium">NickName</label>
          <input type="text" id="nickname" name="nickname" placeholder="John Doe"
            className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
            onChange={(e) => {
              setNickName(e.target.value)
            }}/>
        </div>
        <div className="mb-4">
          <label htmlFor="email" className="block text-gray-700 text-sm font-medium">Email</label>
          <input type="email" id="email" name="email" placeholder="johndoe@example.com"
            className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
            onChange={(e) => {
              setEmail(e.target.value)
            }}/>
        </div>
        <div className="mb-4">
          <label htmlFor="password" className="block text-gray-700 text-sm font-medium">Password</label>
          <input type="password" id="password" name="password" placeholder="********"
            className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
            onChange={(e) => {
              setPassword(e.target.value)
            }}/>
        </div>
        <div className="mb-6">
          <label htmlFor="confirm-password" className="block text-gray-700 text-sm font-medium">Confirm
              Password</label>
          <input type="password" id="confirm-password" name="confirm-password" placeholder="********"
            className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
            onChange={(e) => {
              setConfirmPassword(e.target.value)
            }}/>
        </div>
        <div className="flex items-center justify-between">
          <button type="submit"
            className="bg-indigo-500 text-white py-2 px-4 rounded-md hover:bg-indigo-600 focus:outline-none focus:ring focus:ring-indigo-200 focus:ring-opacity-50">Sign
              Up
          </button>
          <a className="text-indigo-600 hover:text-indigo-700" onClick={switchToLogin}>Already have an account? Log in</a>
        </div>
      </form>
    </Modal>
  )
}

const isStrongPassword = (password: string): boolean => {
  if (password.length < 8) {
    return false
  }

  const uppercaseRegex = /[A-Z]/
  const lowercaseRegex = /[a-z]/
  if (!uppercaseRegex.test(password) || !lowercaseRegex.test(password)) {
    return false
  }

  const digitRegex = /\d/
  return digitRegex.test(password)
}
