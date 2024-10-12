import i18n from '@/i18n'

export const getTranslation = ({ parent, child }: {parent: string; child: string}) => {
  return i18n.global.t(`${parent}.${child?.toLowerCase()}`)
}