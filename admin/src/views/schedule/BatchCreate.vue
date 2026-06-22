<template>
  <div>
    <el-card>
      <template #header>批量排班</template>
      <el-form :model="form" label-width="100px" style="max-width: 600px">
        <el-form-item label="医生">
          <el-select v-model="form.doctorId" placeholder="选择医生" filterable>
            <el-option v-for="d in doctors" :key="d.id" :label="d.name" :value="d.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="form.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="时段">
          <el-checkbox-group v-model="form.periods">
            <el-checkbox value="morning">上午</el-checkbox>
            <el-checkbox value="afternoon">下午</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="号源数">
          <el-input-number v-model="form.totalSlots" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="排除周末">
          <el-switch v-model="form.excludeWeekends" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">生成排班</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { batchCreateSchedule } from '../../api/schedule'
import { getDoctors } from '../../api/doctor'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

const router = useRouter()
const doctors = ref([])
const loading = ref(false)
const form = ref({
  doctorId: null,
  dateRange: [],
  periods: ['morning', 'afternoon'],
  totalSlots: 20,
  excludeWeekends: true
})

onMounted(async () => {
  doctors.value = await getDoctors({}).catch(() => [])
})

async function handleSubmit() {
  if (!form.value.doctorId || !form.value.dateRange?.length) {
    ElMessage.warning('请填写完整信息')
    return
  }
  loading.value = true
  try {
    await batchCreateSchedule({
      doctorId: form.value.doctorId,
      startDate: form.value.dateRange[0],
      endDate: form.value.dateRange[1],
      periods: form.value.periods,
      totalSlots: form.value.totalSlots,
      excludeWeekends: form.value.excludeWeekends
    })
    ElMessage.success('批量排班成功')
    router.push('/schedules')
  } catch (e) {
    ElMessage.error(e.message || '批量排班失败')
  } finally {
    loading.value = false
  }
}
</script>
