const api = require('../../utils/api')

Page({
  data: {
    stats: {
      total: 0,
      paid: 0,
      completed: 0,
      cancelled: 0
    },
    appointments: []
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
    const today = new Date().toISOString().split('T')[0]
    try {
      const [stats, appointments] = await Promise.all([
        api.getAppointmentStats({ date: today }),
        api.getAppointments({ date: today })
      ])
      this.setData({
        stats: stats || this.data.stats,
        appointments: appointments || []
      })
    } catch (err) {
      console.error('加载数据失败:', err)
    }
  }
})
