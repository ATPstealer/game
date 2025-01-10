import { createRouter, createWebHistory } from 'vue-router/auto'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL)
})

router.beforeEach((to, from, next) => {
  if (to.path === '/login') {
    next()
  } else {
    localStorage.setItem('prev', to.path)
  }

  if (localStorage.getItem('invalid') === null) {
    localStorage.setItem('invalid', 'true')
  }

  if (localStorage.getItem('invalid') === 'true') {
    next({ name: 'Login' })
  }

  next()
})

export default router
