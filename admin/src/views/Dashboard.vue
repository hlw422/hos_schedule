<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>今日预约</template>
          <div class="stat-value">{{ stats.total }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>已支付</template>
          <div class="stat-value" style="color: #52C41A">{{ stats.paid }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>已完成</template>
          <div class="stat-value" style="color: #1677FF">{{ stats.completed }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>已取消</template>
          <div class="stat-value" style="color: #FF4D4F">{{ stats.cancelled }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAppointmentStats } from '../api/appointment'

const stats = ref({ total: 0, paid: 0, completed: 0, cancelled: 0 })

onMounted(async () => {
  try {
    const data = await getAppointmentStats()
    stats.value = data
  } catch (e) {
    console.error(e)
  }
})
</script>

<style scoped>
.stat-value {
  font-size: 36px;
  font-weight: bold;
  text-align: center;
  color: #333;
}
</style>
