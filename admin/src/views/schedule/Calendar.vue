<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <div style="display: flex; gap: 12px; align-items: center">
            <el-select v-model="filters.hospitalId" placeholder="医院" clearable style="width: 180px" @change="fetchData">
              <el-option v-for="h in hospitals" :key="h.id" :label="h.name" :value="h.id" />
            </el-select>
            <el-select v-model="filters.departmentId" placeholder="科室" clearable style="width: 150px" @change="fetchData">
              <el-option v-for="d in departments" :key="d.id" :label="d.name" :value="d.id" />
            </el-select>
            <el-date-picker v-model="filters.date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" @change="fetchData" />
          </div>
          <div>
            <el-button type="primary" @click="openAdd">新增排班</el-button>
            <el-button @click="$router.push('/schedules/batch')">批量排班</el-button>
          </div>
        </div>
      </template>
      <el-table :data="schedules" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="doctorName" label="医生" width="100" />
        <el-table-column prop="departmentName" label="科室" width="120" />
        <el-table-column prop="date" label="日期" width="120" />
        <el-table-column label="时段" width="100">
          <template #default="{ row }">
            <el-tag :type="row.period === 'morning' ? '' : 'warning'">
              {{ row.period === 'morning' ? '上午' : '下午' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="totalSlots" label="总号源" width="80" />
        <el-table-column prop="bookedSlots" label="已约" width="80" />
        <el-table-column label="剩余" width="80">
          <template #default="{ row }">
            {{ row.totalSlots - row.bookedSlots }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
          </template>
        </el-table-column>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑排班' : '新增排班'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="医生">
          <el-select v-model="form.doctorId" placeholder="选择医生" filterable>
            <el-option v-for="d in allDoctors" :key="d.id" :label="d.name" :value="d.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="form.date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="时段">
          <el-radio-group v-model="form.period">
            <el-radio value="morning">上午</el-radio>
            <el-radio value="afternoon">下午</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="号源数">
          <el-input-number v-model="form.totalSlots" :min="1" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getSchedules, createSchedule, updateSchedule } from '../../api/schedule'
import { getHospitals } from '../../api/hospital'
import { getDepartments } from '../../api/department'
import { getDoctors } from '../../api/doctor'
import { ElMessage } from 'element-plus'

const schedules = ref([])
const hospitals = ref([])
const departments = ref([])
const allDoctors = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filters = ref({ hospitalId: null, departmentId: null, date: '' })
const form = ref({ id: null, doctorId: null, date: '', period: 'morning', totalSlots: 20 })

onMounted(async () => {
  const [h, d, doc] = await Promise.all([
    getHospitals().catch(() => []),
    getDepartments({}).catch(() => []),
    getDoctors({}).catch(() => [])
  ])
  hospitals.value = h
  departments.value = d
  allDoctors.value = doc
  fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const params = { page: page.value, pageSize: pageSize.value, ...filters.value }
    const res = await getSchedules(params)
    schedules.value = res.list || res
    total.value = res.total || schedules.value.length
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

function openAdd() {
  form.value = { id: null, doctorId: null, date: '', period: 'morning', totalSlots: 20 }
  dialogVisible.value = true
}

function openEdit(row) {
  form.value = { ...row }
  dialogVisible.value = true
}

async function handleSave() {
  try {
    if (form.value.id) {
      await updateSchedule(form.value.id, form.value)
    } else {
      await createSchedule(form.value)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  }
}
</script>
