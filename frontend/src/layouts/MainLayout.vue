<template>
  <q-layout view="lHh lpr lFf" class="app-layout">

    <!-- HEADER -->
    <q-header elevated class="app-header">
      <q-toolbar class="q-px-md q-py-sm">

        <!-- Logo + Título -->
        <div class="row items-center no-wrap">
          <q-avatar size="42px" class="q-mr-md shadow-2">
            <img src="/brand/logo.png" />
          </q-avatar>

          <div class="column">
            <div class="text-weight-bold text-subtitle1">
              IEQ Maria Goretti
            </div>
            <div class="text-caption text-grey-3">
              Escalas da Mídia
            </div>
          </div>
        </div>

        <q-space />

        <!-- Ações -->
        <div class="row items-center q-gutter-sm">

          <!-- Público/Admin -->
          <q-btn
            v-if="!isAdminRoute"
            outline
            color="white"
            class="btn-ghost"
            icon="admin_panel_settings"
            label="Área Admin"
            to="/admin/login"
            no-caps
          />

          <q-btn
            v-else
            outline
            color="white"
            class="btn-ghost"
            icon="public"
            label="Ver público"
            to="/escalas"
            no-caps
          />

          <!-- Sair (só quando estiver logado e no admin) -->
          <q-btn
            v-if="isAdminRoute && isLogged"
            flat
            color="white"
            icon="logout"
            label="Sair"
            no-caps
            @click="logout"
          />
        </div>
      </q-toolbar>
    </q-header>

    <!-- CONTEÚDO -->
    <q-page-container class="app-container">
      <div class="container">
        <router-view />
      </div>
    </q-page-container>

  </q-layout>
</template>

<script>
import { useAuthStore } from 'stores/auth'

export default {
  name: 'MainLayout',

  computed: {
    isAdminRoute () {
      return this.$route.path.startsWith('/admin')
    },

    isLogged () {
      const auth = useAuthStore()
      return auth.isLogged
    }
  },

  methods: {
    logout () {
      const auth = useAuthStore()
      auth.logout()
      this.$router.push('/admin/login')
      this.$q.notify({ type: 'positive', message: 'Sessão encerrada' })
    }
  }
}
</script>

<style scoped>
/* Background do app */
.app-layout {
  background: #f7f7fb;
}

/* Header moderno (gradiente + borda) */
.app-header {
  background: linear-gradient(90deg, #610659 0%, #3a054e 100%);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

/* Container central */
.app-container {
  padding: 18px;
}

/* Limita largura em telas grandes */
.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

/* Botão outline branco mais bonito */
.btn-ghost {
  border-color: rgba(255, 255, 255, 0.55) !important;
}
</style>
