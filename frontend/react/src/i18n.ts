import i18next from 'i18next'
import { initReactI18next } from 'react-i18next'
import en from './locales/en.json'
import ru from './locales/ru.json'

export const defaultNS = 'locale'

// eslint-disable-next-line import/no-named-as-default-member
i18next.use(initReactI18next).init({
  debug: false,
  fallbackLng: localStorage.getItem('lang') || 'en',
  defaultNS,
  resources: {
    en: {
      locale: en
    },
    ru: {
      locale: ru
    }
  }
})

export default i18next