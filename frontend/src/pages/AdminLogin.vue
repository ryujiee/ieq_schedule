<template>
  <q-page class="login-page">
    <q-card class="login-card">

      <q-card-section class="login-header">
        <q-avatar size="72px" class="q-mb-md">
          <img src="/brand/logo.png" />
        </q-avatar>

        <div class="text-h6 text-weight-bold">Área Administrativa</div>
        <div class="text-caption text-grey-7">Acesso restrito à equipe</div>
      </q-card-section>

      <q-separator />

      <q-card-section class="q-pt-lg">
        <q-form @submit.prevent="login" class="q-gutter-sm">

          <q-input v-model="email" type="email" label="Email" outlined dense clearable autofocus :disable="loading"
            class="full-width">
            <template #prepend>
              <q-icon name="mail" />
            </template>
          </q-input>

          <q-input v-model="password" :type="showPassword ? 'text' : 'password'" label="Senha" outlined dense
            :disable="loading" class="full-width">
            <template #prepend>
              <q-icon name="lock" />
            </template>

            <template #append>
              <q-btn flat round dense :icon="showPassword ? 'visibility_off' : 'visibility'"
                @click="showPassword = !showPassword" :tabindex="-1" />
            </template>
          </q-input>

          <q-btn type="submit" label="Entrar" color="primary" class="full-width q-mt-md" :loading="loading" unelevated
            no-caps size="lg" />
        </q-form>
      </q-card-section>

    </q-card>
  </q-page>
</template>

<script>
import { useAuthStore } from 'stores/auth'
import { api } from 'boot/api'

export default {
  name: 'AdminLogin',

  data() {
    return {
      email: '',
      password: '',
      showPassword: false,
      loading: false
    }
  },

  methods: {
    async login() {
      if (!this.email || !this.password) {
        this.$q.notify({ type: 'warning', message: 'Informe email e senha' })
        return
      }

      this.loading = true
      const auth = useAuthStore()

      try {
        await auth.login(this.email, this.password)
        this.$router.push('/admin')
      } catch (err) {
        console.error('Erro ao logar:', err)
        this.$q.notify({ type: 'negative', message: 'Email ou senha inválidos' })
      } finally {
        this.loading = false
      }
    }
  },
  async mounted() {
    // se tem token, tenta validar
    const auth = useAuthStore()
    const token = localStorage.getItem('token')

    if (token && auth.isLogged) {
      try {
        // bate no /health (ou qualquer endpoint admin simples)
        await api.get('/admin/functions')
        this.$router.replace('/admin')
        return
      } catch (e) {
        console.warn('Token inválido ou expirado:', e)
        // token inválido/expirado: limpa
        auth.logout()
      }
    }
  }
}
</script>

<style scoped>
/* Centralização perfeita */
.login-page {
  min-height: calc(100vh - 50px);
  /* header */
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: #f7f7fb;
}

/* Card alinhado e consistente */
.login-card {
  width: 420px;
  max-width: 92vw;
  border-radius: 16px;
  overflow: hidden;
}

/* Header central */
.login-header {
  padding: 28px 28px 18px 28px;
  text-align: center;
}
</style>
