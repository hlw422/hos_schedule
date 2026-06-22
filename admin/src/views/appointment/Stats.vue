<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>预约趋势</span>
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="YYYY-MM-DD"
                @change="fetchStats"
              />
            </div>
          </template>
          <div ref="chartRef" style="height: 400px"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>状态分布</template>
          <div v-for="item in statusStats" :key="item.status" style="padding: 12px 0; display: flex; justify-content: space-between; border-bottom: 1px solid #eee">
            <span>{{ statusLabel(item.status) }}</span>
            <el-tag :type="statusTagType(item.status)">{{ item.count }}</el-tag>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { getAppointmentStats } from '../../api/appointment'
import { ElMessage } from 'element-plus'

const chartRef = ref(null)
const dateRange = ref([])
const statusStats = ref([])

onMounted(fetchStats)

async function fetchStats() {
  try {
    const params = {}
    if (dateRange.value?.length === 2) {
      params.startDate = dateRange.value[0]
      params.endDate = dateRange.value[1]
    }
    const data = await getAppointmentStats(params)
    statusStats.value = data.statusStats || []
    await nextTick()
    renderChart(data.trend || [])
  } catch (e) {
    ElMessage.error('加载统计数据失败')
  }
}

function renderChart(trend) {
  if (!chartRef.value) return
  const el = chartRef.value
  el.innerHTML = ''
  if (!trend.length) {
    el.innerHTML = '<div style="text-align:center;padding-top:180px;color:#999">暂无数据</div>'
    return
  }
  const maxVal = Math.max(...trend.map(t => t.count), 1)
  const barWidth = Math.min(40, (el.clientWidth - 60) / trend.length - 8)
  let html = '<div style="display:flex;align-items:flex-end;gap:4px;height:360px;padding:20px 10px 40px;position:relative">'
  trend.forEach(t => {
    const h = (t.count / maxVal) * 300
    html += `<div style="flex:1;display:flex;flex-direction:column;align-items:center;justify-content:flex-end;height:100%">
      <span style="font-size:12px;margin-bottom:4px">${t.count}</span>
      <div style="width:${barWidth}px;height:${h}px;background:#1677FF;border-radius:4px 4px 0 0"></div>
      <span style="font-size:11px;margin-top:4px;white-space:nowrap">${t.date?.slice(5) || ''}</span>
    </div>`
  })
  html += '</div>'
  el.innerHTML = html
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
