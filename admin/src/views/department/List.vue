<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>科室管理</span>
          <el-button type="primary" @click="openAdd">新增科室</el-button>
        </div>
      </template>
      <el-table :data="departments" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="科室名称" />
        <el-table-column prop="hospitalName" label="所属医院" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑科室' : '新增科室'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="所属医院">
          <el-select v-model="form.hospitalId" placeholder="选择医院">
            <el-option v-for="h in hospitals" :key="h.id" :label="h.name" :value="h.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
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
import { getDepartments, createDepartment, updateDepartment, deleteDepartment } from '../../api/department'
import { getHospitals } from '../../api/hospital'
import { ElMessage, ElMessageBox } from 'element-plus'

const departments = ref([])
const hospitals = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const form = ref({ id: null, name: '', description: '', hospitalId: null })

onMounted(async () => {
  hospitals.value = await getHospitals().catch(() => [])
  fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await getDepartments({ page: page.value, pageSize: pageSize.value })
    departments.value = res.list || res
    total.value = res.total || departments.value.length
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

function openAdd() {
  form.value = { id: null, name: '', description: '', hospitalId: hospitals.value[0]?.id }
  dialogVisible.value = true
}

function openEdit(row) {
  form.value = { ...row }
  dialogVisible.value = true
}

async function handleSave() {
  try {
    if (form.value.id) {
      await updateDepartment(form.value.id, form.value)
    } else {
      await createDepartment(form.value)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm('确定删除该科室？', '提示', { type: 'warning' })
    await deleteDepartment(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || '删除失败')
  }
}
</script>
