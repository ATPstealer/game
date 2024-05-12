import { computed } from 'vue'
import i18n from '@/i18n'

export const timeValues = computed(() => [
  {
    value: 60 * 1000000000,
    label: `${i18n.global.t('common.minute', 1)}`
  },
  {
    value: 3600 * 1000000000,
    label: `${i18n.global.t('common.hour', 1)}`
  },
  {
    value: 7200 * 1000000000,
    label: `${i18n.global.t('common.hour', 2)}`
  },
  {
    value: 18000 * 1000000000,
    label: `${i18n.global.t('common.hour', 5)}`
  },
  {
    value: 86400 * 1000000000,
    label: `${i18n.global.t('common.hour', 24)}`
  },
  {
    value: 604800 * 1000000000,
    label: `${i18n.global.t('common.day', 7)}`
  },
  {
    value: 2592000 * 1000000000,
    label: `${i18n.global.t('common.day', 30)}`
  }
])