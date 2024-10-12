import '@/assets/main.css'
import '@/assets/tailwind.css'
import 'primevue/resources/themes/lara-light-teal/theme.css'
import 'primeicons/primeicons.css'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ConfirmationService from 'primevue/confirmationservice'
import ToastService from 'primevue/toastservice'
import Tooltip from 'primevue/tooltip'
import { createApp } from 'vue'
import App from '@/App.vue'
import i18n from '@/i18n'
import router from '@/router'

const vueQueryPluginOptions = {
  queryClientConfig: {
    defaultOptions: {
      queries: {
        select: (data: {data: any}) => data.data
      }
    }
  }
}

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)
app.use(PrimeVue)
app.use(ConfirmationService)
app.use(ToastService)
app.use(VueQueryPlugin, vueQueryPluginOptions)
app.directive('tooltip', Tooltip)

app.mount('#app')
