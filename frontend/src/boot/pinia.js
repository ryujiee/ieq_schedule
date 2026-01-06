import { boot } from 'quasar/wrappers'
import { createPinia } from 'pinia'
import { useAuthStore } from 'stores/auth'

export default boot(({ app }) => {
  const pinia = createPinia()
  app.use(pinia)

  // inicializa token no axios
  const auth = useAuthStore(pinia)
  auth.init()
})
