import isEmpty from 'lodash/isEmpty'
import { Button } from 'primereact/button'
import { Dropdown, DropdownChangeEvent } from 'primereact/dropdown'
import { Menu } from 'primereact/menu'
// eslint-disable-next-line import/no-unresolved
import { MenuItem, MenuItemOptions } from 'primereact/menuitem'
import React, { useEffect, useRef, useState } from 'react'
import { useCookies } from 'react-cookie'
import { useTranslation } from 'react-i18next'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useRecoilState } from 'recoil'
import { useGetEntity } from '../../hooks/getData'
import { logout } from '../../hooks/logout'
import { formatNumber } from '../../include/formatNumber'
import { IUser } from '../../models'
import { nickNameState } from '../../state'
import { Language, userPages, worldPages, languages } from './constants'
import './styles.scss'

const Header = ({ showLogin, showRegister }: HeaderProps) => {
  const [cookies, setCookie] = useCookies(['lang'])
  const defaultLanguage =  cookies.lang || 'en'
  const [language, setLanguage] = useState<Language>(defaultLanguage)
  const [currentPage, setCurrentPage] = useState<string>(useLocation().pathname)
  const user = useGetEntity<IUser>('/user/data')
  const userPagesRef = useRef<any>(null)
  const worldPagesRef = useRef<any>(null)
  const { i18n, t } = useTranslation()
  const navigate = useNavigate()
  const [, setNick] = useRecoilState(nickNameState)

  const handleChangeLanguage = (event: DropdownChangeEvent) => {
    setLanguage(event.target.value as Language)
    setCookie('lang', event.target.value)
    i18n.changeLanguage(event.target.value)
  }

  useEffect(() => {
    i18n.changeLanguage(defaultLanguage)
  },[defaultLanguage, i18n])

  useEffect(() => {
    if (user) {
      const { nickName } = user
      setNick(nickName)
    }
  }, [user])

  const logOut = () => {
    logout()
  }

  const userItems = userPages.map(page => {
    return {
      label: t(`${page}.title`),
      command: () => {
        navigate(`/${page}`)
        setCurrentPage(`/${page}`)
      },
      template: (item: MenuItem, options: MenuItemOptions) => {
        return (
          <a className={'link' + (currentPage === `/${page}` ? ' font-bold': '')} onClick={options.onClick}>
            <span>{item.label}</span>
          </a>
        )
      }
    }
  })

  userItems.unshift({
    label: t('main.title'),
    command: () => {
      navigate('/')
      setCurrentPage('/')
    },
    template: (item: MenuItem, options: MenuItemOptions) => {
      return (
        <a className={'link' + (currentPage === '/' ? ' font-bold': '')} onClick={options.onClick}>
          <span>{item.label}</span>
        </a>
      )
    }
  })

  const worldItems = worldPages.map(page => {
    return {
      label: t(`${page}.title`),
      command: () => {
        navigate(`/${page}`)
        setCurrentPage(`/${page}`)
      },
      template: (item: MenuItem, options: MenuItemOptions) => {
        return (
          <a className={'link' + (currentPage === `/${page}` ? ' font-bold': '')} onClick={options.onClick}>
            <span>{item.label}</span>
          </a>
        )
      }
    }
  })

  worldItems.unshift({
    label: t('map.title'),
    command: () => {
      navigate('/map')
      setCurrentPage('/map')
    },
    template: (item: MenuItem, options: MenuItemOptions) => {
      return (
        <a className={'link' + (currentPage === '/map' ? ' font-bold': '')} onClick={options.onClick}>
          <span>{item.label}</span>
        </a>
      )
    }
  })

  return (
    <div className="bg-gray-500 h-8 py-8 px-4 flex items-center justify-between">
      <div className="flex gap-2">
        <Menu
          model={userItems}
          popup
          ref={userPagesRef}
          id='user-pages'
        />
        <Button
          label={t('menu.business')}
          onClick={
            (event) => userPagesRef.current?.toggle(event)
          }
          aria-controls="user-pages" aria-haspopup
          text
          className="font-semibold text-white focus:shadow-none"
        />
        <Menu
          model={worldItems}
          popup
          ref={worldPagesRef}
          id='world-pages'
        />
        <Button
          label={t('menu.world')}
          onClick={
            (event) => worldPagesRef.current?.toggle(event)
          }
          aria-controls="world-pages" aria-haspopup
          text
          className="font-semibold text-white focus:shadow-none"
        />
      </div>
      <div className="flex gap-4 items-center">
        <div className="flex font-bold justify-start gap-4 ml-auto text-white">
          {
            !isEmpty(user) &&
              <>
                <Link to="/" className="hover:text-gray-950">{user.nickName}</Link>
                <span> {user.money && formatNumber(user.money)}$ </span>
                <button className="hover:text-gray-950" onClick={logOut}>Log out</button>
              </>
          }
        </div>
        <div className="flex gap-4 ml-4 text-white">
          {
            isEmpty(user) &&
              <>
                <button className="mr-2" onClick={() => {showLogin(true)}}>Login</button>
                <button className="mr-2" onClick={() => {showRegister(true)}}>Sing up</button>
              </>
          }
        </div>
        <Dropdown
          value={language}
          options={languages}
          optionLabel='label'
          optionValue='key'
          onChange={handleChangeLanguage}
          className="lang bg-gray-500 border-gray-500 !text-white"
        />
      </div>
    </div>
  )
}

type HeaderProps = {
  showLogin: (show: boolean) => void;
  showRegister: (show: boolean) => void;
}

export default Header