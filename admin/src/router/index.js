import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '../layouts/AdminLayout.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/',
    component: AdminLayout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '数据看板', icon: 'DataBoard' }
      },
      {
        path: 'hospitals',
        name: 'Hospitals',
        component: () => import('../views/hospital/List.vue'),
        meta: { title: '医院管理', icon: 'Hospital' }
      },
      {
        path: 'departments',
        name: 'Departments',
        component: () => import('../views/department/List.vue'),
        meta: { title: '科室管理', icon: 'Grid' }
      },
      {
        path: 'doctors',
        name: 'Doctors',
        component: () => import('../views/doctor/List.vue'),
        meta: { title: '医生管理', icon: 'User' }
      },
      {
        path: 'schedules',
        name: 'Schedules',
        component: () => import('../views/schedule/Calendar.vue'),
        meta: { title: '排班管理', icon: 'Calendar' }
      },
      {
        path: 'schedules/batch',
        name: 'BatchSchedule',
        component: () => import('../views/schedule/BatchCreate.vue'),
        meta: { title: '批量排班', icon: 'Document' }
      },
      {
        path: 'appointments',
        name: 'Appointments',
        component: () => import('../views/appointment/List.vue'),
        meta: { title: '预约管理', icon: 'List' }
      },
      {
        path: 'appointments/stats',
        name: 'AppointmentStats',
        component: () => import('../views/appointment/Stats.vue'),
        meta: { title: '预约统计', icon: 'TrendCharts' }
      }
    ]
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
