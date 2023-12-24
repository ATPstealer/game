import { defaultNS } from '../i18n'
import en from '../locales/en.json'
import ru from '../locales/ru.json'

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: typeof defaultNS;
    resources: {
      en: typeof en;
      ru: typeof ru;
    };
  }
}