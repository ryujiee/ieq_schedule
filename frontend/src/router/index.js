import { route } from 'quasar/wrappers'
import { createRouter, createWebHistory } from 'vue-router'
import routes from './routes'

export default route(function () {
  const Router = createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior: () => ({ left: 0, top: 0 })
  })

  Router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')

    // 🔁 Sempre redireciona "/" para "/escalas"
    if (to.path === '/') {
      return next('/escalas')
    }

    // 🔒 Rotas admin protegidas
    if (to.path.startsWith('/admin')) {
      // libera login
      if (to.path === '/admin/login') {
        return next()
      }

      // se não tiver token → login
      if (!token) {
        return next('/admin/login')
      }
    }

    next()
  })

  return Router
})
