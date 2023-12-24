import React from 'react'
import { useTranslation } from 'react-i18next'
import { useGetEntity } from '../hooks/getData'
import { Characteristics, IUser } from '../models'

const CharacteristicsPage = () => {
  const userData = useGetEntity<IUser>('/user/data')
  const { t } = useTranslation()

  const characteristics: Characteristics = {
    memory: userData?.memory,
    intelligence: userData?.intelligence,
    attention: userData?.attention,
    wits: userData?.wits,
    multitasking: userData?.multitasking,
    management: userData?.management,
    planning: userData?.planning
  }

  return (
    <div className="p-4 mb-4 md:mb-0 shadow-xl">
      <h2 className="font-bold mb-2 text-2xl">{t('user.characteristics')}</h2>
      <div className="grid gap-1">
        {Object.keys(characteristics).map(item => item && <p key={item}>{t(`user.${item}`)}: {characteristics[item as keyof Characteristics]}</p>)}
      </div>
    </div>
  )
}

export default CharacteristicsPage