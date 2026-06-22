<template>
  <div>
    <el-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>院区管理</span>
          <el-button type="primary" @click="openAdd">新增院区</el-button>
        </div>
      </template>
      <el-table :data="campuses" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="院区名称" />
        <el-table-column prop="address" label="地址" />
        <el-table-column prop="phone" label="电话" width="140" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新增院区" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="所属医院">
          <el-select v-model="form.hospitalId" placeholder="选择医院">
            <el-option v-for="h in hospitals" :key="h.id" :label="h.name" :value="h.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="form.name" />
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
import { getHospitals, getCampuses, addCampus } from '../../api/hospital'
import { ElMessage } from 'element-plus'

const campuses = ref([])
const hospitals = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const form = ref({ hospitalId: null, name: '', address: '', phone: '' })

onMounted(async () => {
  loading.value = true
  try {
    hospitals.value = await getHospitals()
    if (hospitals.value.length) {
      campuses.value = await getCampuses(hospitals.value[0].id)
    }
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
})

function openAdd() {
  form.value = { hospitalId: hospitals.value[0]?.id, name: '', address: '', phone: '' }
  dialogVisible.value = true
}

async function handleSave() {
  try {
    await addCampus(form.value)
    ElMessage.success('添加成功')
    dialogVisible.value = false
    if (form.value.hospitalId) {
      campuses.value = await getCampuses(form.value.hospitalId)
    }
  } catch (e) {
    ElMessage.error(e.message || '添加失败')
  }
}
</script>
