<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>预约管理</span>
          <div style="display: flex; gap: 12px">
            <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px" @change="fetchData">
              <el-option label="待支付" value="pending" />
              <el-option label="已支付" value="paid" />
              <el-option label="已完成" value="completed" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
            <el-date-picker v-model="filters.date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" @change="fetchData" />
          </div>
        </div>
      </template>
      <el-table :data="appointments" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="patientName" label="患者" width="100" />
        <el-table-column prop="doctorName" label="医生" width="100" />
        <el-table-column prop="departmentName" label="科室" width="120" />
        <el-table-column prop="date" label="日期" width="120" />
        <el-table-column label="时段" width="80">
          <template #default="{ row }">
            {{ row.period === 'morning' ? '上午' : '下午' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" />
      </el-table>
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end"
        @current-change="fetchData"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAppointments } from '../../api/appointment'
import { ElMessage } from 'element-plus'

const appointments = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filters = ref({ status: '', date: '' })

onMounted(fetchData)

async function fetchData() {
  loading.value = true
  try {
    const params = { page: page.value, pageSize: pageSize.value, ...filters.value }
    const res = await getAppointments(params)
    appointments.value = res.list || res
    total.value = res.total || appointments.value.length
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

function statusLabel(s) {
  const map = { pending: '待支付', paid: '已支付', completed: '已完成', cancelled: '已取消' }
  return map[s] || s
}

function statusTagType(s) {
  const map = { pending: 'warning', paid: '', completed: 'success', cancelled: 'info' }
  return map[s] || ''
}
</script>
