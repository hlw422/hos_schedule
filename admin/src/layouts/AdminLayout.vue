<template>
  <el-container style="height: 100vh">
    <el-aside width="220px" style="background: #001529">
      <div class="logo">
        <h2 style="color: #fff; text-align: center; padding: 20px 0; margin: 0;">医院预约系统</h2>
      </div>
      <el-menu
        :default-active="$route.path"
        background-color="#001529"
        text-color="#ffffffa6"
        active-text-color="#ffffff"
        router
      >
        <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.title }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="background: #fff; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #eee">
        <span style="font-size: 18px; font-weight: bold">{{ $route.meta.title }}</span>
        <el-dropdown @command="handleCommand">
          <span style="cursor: pointer">
            管理员 <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main style="background: #f5f5f5; padding: 20px">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const menuItems = [
  { path: '/dashboard', title: '数据看板', icon: 'DataBoard' },
  { path: '/hospitals', title: '医院管理', icon: 'Hospital' },
  { path: '/departments', title: '科室管理', icon: 'Grid' },
  { path: '/doctors', title: '医生管理', icon: 'User' },
  { path: '/schedules', title: '排班管理', icon: 'Calendar' },
  { path: '/appointments', title: '预约管理', icon: 'List' },
  { path: '/appointments/stats', title: '预约统计', icon: 'TrendCharts' }
]

function handleCommand(cmd) {
  if (cmd === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>
