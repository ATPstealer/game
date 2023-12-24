import axios, { AxiosError } from 'axios'
import sha256 from 'crypto-js/sha256'
import React, { useEffect, useState } from 'react'
import { useCookies } from 'react-cookie'
import { useNavigate } from 'react-router-dom'
import config from '../config'
import { IResponseEntity, IToken, IUser } from '../models'
import { Error } from './Error'
import { Modal } from './Modal'
import { Success } from './Success'

interface LoginProps {
  onClose: () => void;
  switchToRegister: () => void;
}

const userData: IUser = {
  nickName: '',
  email: '',
  password: '',
  money: 0
}

export const Login = ({ onClose, switchToRegister }: LoginProps) => {
  const [nickName, setNickName] = useState('') // TODO: как сделать чтобы призакрытии модалки данные жили
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [, setCookie] = useCookies(['secureToken'])
  const navigate = useNavigate()
  const [isLogin, setIsLogin] = useState(false)

  const submitHandler = async (event: React.FormEvent) => {
    event.preventDefault()
    setError('')
    setMessage('')

    userData.nickName = nickName
    userData.password = sha256(password).toString()
    try {
      const response = await axios.post<IResponseEntity<IToken>>(config.apiBaseUrl + '/user/login', userData)
      if (response.data.status === 'success') {
        const ttl = Number(response.data.data?.ttl)
        const date = new Date()
        date.setTime(date.getTime() + ttl*1000)
        setCookie('secureToken', response.data.data?.token, { expires: date, domain:'.' + config.domain })
        setMessage('You are logged in. You will be redirected')
        setIsLogin(true)
      } else {
        setError(response.data.text)
      }
    } catch
    (e: unknown) {
      const error = e as AxiosError
      setError(error.message)
    }
  }

  useEffect(() => {
    if (isLogin) {
      setTimeout(() => {
        navigate('/')
        onClose()
        window.location.reload()
      }, 2000)
    }
  })

  return (
    <Modal title="Login" onClose={onClose}>
      <Error error={error}/>
      <Success message={message}/>
      <form onSubmit={submitHandler}>
        <div className="mb-4">
          <label htmlFor="nickname" className="block text-gray-700 text-sm font-medium">nickName</label>
          <input type="text" id="nickname" name="nickname" placeholder="John Doe"
            className="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
            onChange={(e) => {
              setNickName(e.target.value)
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
        <div className="flex items-center justify-between">
          <button type="submit"
            className="bg-indigo-500 text-white py-2 px-4 rounded-md hover:bg-indigo-600 focus:outline-none focus:ring focus:ring-indigo-200 focus:ring-opacity-50">Login
          </button>
          <a className="text-indigo-600 hover:text-indigo-700" onClick={switchToRegister}>Sign up</a>
        </div>
      </form>
    </Modal>
  )
}