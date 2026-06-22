<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>医生管理</span>
          <el-button type="primary" @click="openAdd">新增医生</el-button>
        </div>
      </template>
      <el-table :data="doctors" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="title" label="职称" width="120" />
        <el-table-column prop="departmentName" label="科室" width="150" />
        <el-table-column prop="specialty" label="擅长" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? '在岗' : '停诊' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button
              size="small"
              :type="row.status === 'active' ? 'warning' : 'success'"
              @click="toggleStatus(row)"
            >
              {{ row.status === 'active' ? '停诊' : '启用' }}
            </el-button>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑医生' : '新增医生'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="科室">
          <el-select v-model="form.departmentId" placeholder="选择科室">
            <el-option v-for="d in departments" :key="d.id" :label="d.name" :value="d.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="职称">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="擅长">
          <el-input v-model="form.specialty" type="textarea" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.bio" type="textarea" :rows="3" />
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
import { getDoctors, createDoctor, updateDoctor, updateDoctorStatus } from '../../api/doctor'
import { getDepartments } from '../../api/department'
import { ElMessage } from 'element-plus'

const doctors = ref([])
const departments = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const form = ref({ id: null, name: '', departmentId: null, title: '', specialty: '', bio: '' })

onMounted(async () => {
  departments.value = await getDepartments({}).catch(() => [])
  fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await getDoctors({ page: page.value, pageSize: pageSize.value })
    doctors.value = res.list || res
    total.value = res.total || doctors.value.length
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

function openAdd() {
  form.value = { id: null, name: '', departmentId: null, title: '', specialty: '', bio: '' }
  dialogVisible.value = true
}

function openEdit(row) {
  form.value = { ...row }
  dialogVisible.value = true
}

async function handleSave() {
  try {
    if (form.value.id) {
      await updateDoctor(form.value.id, form.value)
    } else {
      await createDoctor(form.value)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  }
}

async function toggleStatus(row) {
  try {
    const newStatus = row.status === 'active' ? 'inactive' : 'active'
    await updateDoctorStatus(row.id, newStatus)
    ElMessage.success('操作成功')
    fetchData()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}
</script>
