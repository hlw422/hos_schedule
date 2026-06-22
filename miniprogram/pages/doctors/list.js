const api = require('../../utils/api')

Page({
  data: {
    departmentId: '',
    currentFilter: 'all',
    doctors: [],
    loading: false
  },

  onLoad(options) {
    this.setData({ departmentId: options.department_id })
    this.loadDoctors()
  },

  async loadDoctors() {
    this.setData({ loading: true })
    try {
      const params = { department_id: this.data.departmentId }
      if (this.data.currentFilter !== 'all') {
        params.period = this.data.currentFilter
      }
      const doctors = await api.getDoctors(params)
      this.setData({ doctors })
    } catch (e) {
      console.error('Failed to load doctors:', e)
    } finally {
      this.setData({ loading: false })
    }
  },

  setFilter(e) {
    const filter = e.currentTarget.dataset.filter
    this.setData({ currentFilter: filter })
    this.loadDoctors()
  },

  goDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({ url: `/pages/doctors/detail?id=${id}` })
  }
})
