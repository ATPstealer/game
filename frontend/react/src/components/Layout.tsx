import React, { useState } from 'react'
import { Outlet } from 'react-router-dom'
import Header from './Header/Header'
import { Login } from './Login'
import { Register } from './Register'

export const Layout = () => {
  const [registerWindow, setRegisterWindow] = useState(false)
  const [loginWindow, setLoginWindow] = useState(false)

  return (
    <>
      <Header
        showLogin={setLoginWindow}
        showRegister={setRegisterWindow}
      />

      <Outlet/>

      <div className="mt-auto">
        Footer
      </div>
      {registerWindow && <Register onClose={() => setRegisterWindow(false)}
        switchToLogin={() => {
          setRegisterWindow(false)
          setLoginWindow(true)
        }
        }/>} {/* TODO: спросить как вынести нахер*/}
      {loginWindow && <Login onClose={() => setLoginWindow(false)}
        switchToRegister={() => {
          setRegisterWindow(true)
          setLoginWindow(false)
        }
        }/>}

    </>
  )
}
