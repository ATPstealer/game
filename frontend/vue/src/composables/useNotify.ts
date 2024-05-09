import { useToast } from 'primevue/usetoast'

// хак, для показа уведомлений из useAuthFetch
declare global {
  interface Window {
    notify: any;
  }
}

type Status = 'success' | 'error' | 'warn' | 'info'
interface Notify {
  status: Status;
  label: string;
}

export const useNotify = () => {
  const toast = useToast()

  const setNotify = (notify: Notify) => {
    toast.add({ severity: notify.status, summary: notify.label, life: 5000 })
  }

  const setWarning = (label: string) => {
    setNotify({ status: 'warn', label })
  }

  const setError = (label: string) => {
    setNotify({ status: 'error', label })
  }

  const setSuccess = (label: string) => {
    setNotify({ status: 'success', label })
  }

  window.notify = {
    setError
  }

  return {
    setNotify,
    setWarning,
    setError,
    setSuccess
  }
}
