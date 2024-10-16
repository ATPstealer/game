import { createRouter, createWebHistory } from 'vue-router/auto'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL)
})
// TODO: СПРОСИТЬ
router.beforeEach((to, from, next) => {
  localStorage.setItem('prev', from.path)

  if (to.path === '/login') {
    next()
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
