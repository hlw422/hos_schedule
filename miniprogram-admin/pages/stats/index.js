const api = require('../../utils/api')

Page({
  data: {
    trend: [],
    monthStats: {
      total: 0,
      completed: 0,
      cancelled: 0,
      rate: 0
    }
  },

  onLoad() {
    this.loadData()
  },

  onPullDownRefresh() {
    this.loadData().then(() => {
      wx.stopPullDownRefresh()
    })
  },

  async loadData() {
    const now = new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const params = { year, month }

    try {
      const [trend, stats] = await Promise.all([
        api.getStatsTrend(params),
        api.getAppointmentStats(params)
      ])
      this.setData({
        trend: trend || [],
        monthStats: {
          total: stats?.total || 0,
          completed: stats?.completed || 0,
          cancelled: stats?.cancelled || 0,
          rate: stats?.total ? Math.round((stats.completed / stats.total) * 100) : 0
        }
      })
    } catch (err) {
      console.error('加载统计数据失败:', err)
    }
  }
})
