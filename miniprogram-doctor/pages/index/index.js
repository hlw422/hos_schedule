const api = require('../../utils/api')

Page({
  data: {
    today: '',
    appointments: []
  },
  onLoad() {
    this.setData({ today: new Date().toLocaleDateString('zh-CN') })
    this.loadAppointments()
  },
  onShow() {
    this.loadAppointments()
  },
  async loadAppointments() {
    try {
      const data = await api.getTodayAppointments()
      this.setData({ appointments: data || [] })
    } catch (e) {
      console.error(e)
    }
  }
})
