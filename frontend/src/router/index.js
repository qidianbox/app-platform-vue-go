import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue')
  },
  {
    path: '/',
    component: () => import('@/layouts/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue')
      },
      {
        path: 'apps',
        name: 'Apps',
        component: () => import('@/views/app/index.vue')
      },
      {
        path: 'modules',
        name: 'Modules',
        component: () => import('@/views/module/index.vue')
      },
      {
        path: 'system/audit',
        name: 'AuditLog',
        component: () => import('@/views/system/AuditLog.vue'),
        meta: { title: '操作审计日志' }
      }
    ]
  },
  // APP详情页面独立于主布局，拥有自己的顶部导航
  {
    path: '/apps/:appId',
    component: () => import('@/views/app/config/Layout.vue'),
    children: [
      {
        path: 'config',
        name: 'AppConfig',
        component: () => import('@/views/app/config/index.vue')
      },
      {
        path: 'menus',
        name: 'AppMenus',
        component: () => import('@/views/app/config/MenuManagement.vue'),
        meta: { title: '菜单管理' }
      },
      {
        path: 'apis',
        name: 'AppAPIs',
        component: () => import('@/views/app/config/APIManagement.vue'),
        meta: { title: 'API管理' }
      }
    ]
  },
  // 兼容旧路由
  {
    path: '/apps/:id/config',
    redirect: to => ({ path: `/apps/${to.params.id}/config` })
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
