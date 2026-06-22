const api = require('../../utils/api')

Page({
  data: {
    currentStatus: 'all',
    appointments: [],
    loading: false
  },

  onShow() {
    this.loadAppointments()
  },

  async loadAppointments() {
    this.setData({ loading: true })
    try {
      const params = {}
      if (this.data.currentStatus !== 'all') {
        params.status = this.data.currentStatus
      }
      const appointments = await api.getAppointments(params)
      this.setData({ appointments })
    } catch (e) {
      console.error('Failed to load appointments:', e)
    } finally {
      this.setData({ loading: false })
    }
  },

  setStatus(e) {
    const status = e.currentTarget.dataset.status
    this.setData({ currentStatus: status })
    this.loadAppointments()
  },

  goDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({ url: `/pages/appointment/success?id=${id}` })
  },

  cancelAppointment(e) {
    const id = e.currentTarget.dataset.id
    wx.showModal({
      title: '取消预约',
      content: '确定要取消该预约吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            await api.cancelAppointment(id, '用户主动取消')
            wx.showToast({ title: '取消成功', icon: 'success' })
            this.loadAppointments()
          } catch (e) {
            wx.showToast({ title: e || '取消失败', icon: 'none' })
          }
        }
      }
    })
  }
})
