<template>
  <div>
    <el-card>
      <el-table :data="hospitals" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="医院名称" />
        <el-table-column prop="level" label="医院等级" width="120" />
        <el-table-column prop="address" label="地址" />
        <el-table-column prop="phone" label="电话" width="140" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" @click="$router.push(`/hospitals/${row.id}/campuses`)">院区</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="编辑医院" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="等级">
          <el-input v-model="form.level" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.address" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="form.phone" />
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
import { getHospitals, updateHospital } from '../../api/hospital'
import { ElMessage } from 'element-plus'

const hospitals = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const form = ref({ id: null, name: '', level: '', address: '', phone: '' })

onMounted(fetchData)

async function fetchData() {
  loading.value = true
  try {
    hospitals.value = await getHospitals()
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

function openEdit(row) {
  form.value = { ...row }
  dialogVisible.value = true
}

async function handleSave() {
  try {
    await updateHospital(form.value.id, form.value)
    ElMessage.success('保存成功')
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  }
}
</script>
