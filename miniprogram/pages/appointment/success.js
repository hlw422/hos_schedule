const api = require('../../utils/api')

Page({
  data: {
    appointment: {}
  },

  onLoad(options) {
    this.loadAppointment(options.id)
  },

  async loadAppointment(id) {
    try {
      const appointment = await api.getAppointment(id)
      this.setData({ appointment })
    } catch (e) {
      console.error('Failed to load appointment:', e)
    }
  },

  goList() {
    wx.switchTab({ url: '/pages/appointment/list' })
  },

  goHome() {
    wx.switchTab({ url: '/pages/index/index' })
  }
})
