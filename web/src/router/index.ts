import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { getAuthToken } from '@/utils/auth'

// 页面组件懒加载
const Home = () => import('@/views/Home.vue')
const TimeLog = () => import('@/views/TimeLog.vue')
const Tasks = () => import('@/views/Tasks.vue')
const Categories = () => import('@/views/Categories.vue')
const Statistics = () => import('@/views/Statistics.vue')
// @ts-ignore
const Constraints = () => import('@/views/Constraints.vue') as Promise<{ default: any }>
const Metrics = () => import('@/views/Metrics.vue')
const Study = () => import('@/views/Study.vue')
const PasskeyRegister = () => import('@/views/PasskeyRegister.vue')
const PasskeyLogin = () => import('@/views/PasskeyLogin.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: 'Dashboard',
    },
  },
  {
    path: '/timelogs',
    name: 'TimeLog',
    component: TimeLog,
    meta: {
      title: 'Time Logs',
    },
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: Tasks,
    meta: {
      title: 'Task Management',
    },
  },
  {
    path: '/categories',
    name: 'Categories',
    component: Categories,
    meta: {
      title: 'Category Management',
    },
  },
  {
    path: '/statistics',
    name: 'Statistics',
    component: Statistics,
    meta: {
      title: 'Statistics',
    },
  },
  {
    path: '/constraints',
    name: 'Constraints',
    component: Constraints,
    meta: {
      title: '约束',
    },
  },
  {
    path: '/metrics',
    name: 'Metrics',
    component: Metrics,
    meta: {
      title: '指标',
    },
  },
  {
    path: '/study',
    name: 'Study',
    component: Study,
    meta: {
      title: '学习情况',
    },
  },
  {
    path: '/passkey/register',
    name: 'PasskeyRegister',
    component: PasskeyRegister,
    meta: {
      title: 'Passkey Setup',
    },
  },
  {
    path: '/login',
    name: 'PasskeyLogin',
    component: PasskeyLogin,
    meta: {
      title: 'Passkey Login',
    },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
})
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - TimeLog`
  }

  // 认证检查
  const token = getAuthToken()
  const publicPages = ['PasskeyLogin', 'PasskeyRegister']
  const authRequired = !publicPages.includes(to.name as string)

  if (authRequired && !token) {
    next('/login')
  } else if (!authRequired && token && to.name === 'PasskeyLogin') {
    // 如果已登录且访问登录页，重定向到首页
    next('/')
  } else {
    next()
  }
})

export default router
