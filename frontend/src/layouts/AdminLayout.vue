<template>
  <q-layout view="lHh Lpr lFf" class="admin-layout">

    <!-- Header -->
    <q-header elevated class="admin-header">
      <q-toolbar class="q-px-md q-py-sm">
        <q-btn flat round dense icon="menu" class="text-white q-mr-sm" @click="leftDrawerOpen = !leftDrawerOpen" />

        <div class="row items-center no-wrap">
          <q-avatar size="34px" class="q-mr-sm shadow-2">
            <img src="/brand/logo.png" />
          </q-avatar>

          <div class="column">
            <div class="text-weight-bold text-white">IEQ Maria Goretti</div>
            <div class="text-caption text-grey-3">Admin • Escalas da Mídia</div>
          </div>
        </div>

        <q-space />

        <q-btn outline color="white" class="btn-ghost" icon="public" label="Ver público" to="/escalas" no-caps />
      </q-toolbar>
    </q-header>

    <!-- Sidebar -->
    <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="admin-drawer" :width="260">
      <div class="drawer-top q-pa-md">
        <div class="text-subtitle1 text-weight-bold">Painel Admin</div>
        <div class="text-caption text-grey-7">Gerenciamento do sistema</div>
      </div>

      <q-separator />

      <q-list padding class="q-pt-sm">
        <q-item clickable v-ripple :active="isActive('/admin')" active-class="nav-active" to="/admin">
          <q-item-section avatar>
            <q-icon name="calendar_month" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Dashboard</q-item-label>
            <q-item-label caption>Calendário e escalas</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable v-ripple :active="isActive('/admin/funcoes')" active-class="nav-active" to="/admin/funcoes">
          <q-item-section avatar>
            <q-icon name="badge" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Funções</q-item-label>
            <q-item-label caption>Cadastro de funções</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable v-ripple :active="isActive('/admin/membros')" active-class="nav-active" to="/admin/membros">
          <q-item-section avatar>
            <q-icon name="groups" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Membros</q-item-label>
            <q-item-label caption>Cadastro de membros</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable v-ripple to="/admin/whatsapp" active-class="menu-active">
          <q-item-section avatar>
            <q-icon name="chat" />
          </q-item-section>
          <q-item-section>
            <q-item-label>WhatsApp</q-item-label>
            <q-item-label caption>Configuração QRCode</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable v-ripple to="/admin/settings" active-class="menu-active">
          <q-item-section avatar>
            <q-icon name="settings" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Configurações</q-item-label>
            <q-item-label caption>Configuração do Sistema</q-item-label>

          </q-item-section>
        </q-item>

        <q-separator class="q-my-sm" />

        <q-item clickable v-ripple @click="logout">
          <q-item-section avatar>
            <q-icon name="logout" color="negative" />
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-negative">Sair</q-item-label>
            <q-item-label caption>Encerrar sessão</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <!-- Conteúdo -->
    <q-page-container class="admin-container">
      <div class="container">
        <router-view />
      </div>
    </q-page-container>

  </q-layout>
</template>

<script>
import { useAuthStore } from 'stores/auth'

export default {
  name: 'AdminLayout',

  data() {
    return {
      leftDrawerOpen: true
    }
  },

  methods: {
    isActive(path) {
      return this.$route.path === path
    },

    logout() {
      const auth = useAuthStore()
      auth.logout()
      this.$router.push('/admin/login')
      this.$q.notify({ type: 'positive', message: 'Sessão encerrada' })
    }
  }
}
</script>

<style scoped>
.admin-layout {
  background: #f7f7fb;
}

.admin-header {
  background: linear-gradient(90deg, #610659 0%, #3a054e 100%);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.admin-drawer {
  background: #ffffff;
}

.drawer-top {
  background: linear-gradient(180deg, rgba(97, 6, 89, 0.08), rgba(58, 5, 78, 0.02));
}

.admin-container {
  padding: 18px;
}

.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}

.btn-ghost {
  border-color: rgba(255, 255, 255, 0.55) !important;
}

.nav-active {
  background: rgba(97, 6, 89, 0.10) !important;
  border-left: 4px solid #610659;
}
</style>
