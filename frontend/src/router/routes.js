const routes = [
  // Público
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', redirect: '/escalas' },
      { path: 'escalas', component: () => import('pages/PublicSchedule.vue') }
    ]
  },

  // Admin (com sidebar)
  {
    path: '/admin',
    component: () => import('layouts/AdminLayout.vue'), // ou MainLayout com sidebar admin
    children: [
      { path: '', component: () => import('pages/AdminDashboard.vue') },
      { path: 'funcoes', component: () => import('pages/AdminFunctions.vue') },
      { path: 'membros', component: () => import('pages/AdminMembers.vue') },
      { path: 'whatsapp', component: () => import('pages/AdminWhatsApp.vue') },
      { path: 'settings', component: () => import('pages/AdminSettings.vue') }
    ]
  },


  // Login (sem sidebar, pode ficar no MainLayout ou sozinho)
  {
    path: '/admin/login',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/AdminLogin.vue') }
    ]
  },

  // 404
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue')
  }
]

export default routes
