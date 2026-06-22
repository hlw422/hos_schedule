const api = require('../../utils/api')

Page({
  data: {
    hospitalId: '',
    departments: [],
    loading: false
  },

  onLoad(options) {
    this.setData({ hospitalId: options.hospital_id })
    this.loadDepartments()
  },

  async loadDepartments() {
    this.setData({ loading: true })
    try {
      const departments = await api.getDepartments({ hospital_id: this.data.hospitalId })
      this.setData({ departments })
    } catch (e) {
      console.error('Failed to load departments:', e)
    } finally {
      this.setData({ loading: false })
    }
  },

  goDoctors(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({ url: `/pages/doctors/list?department_id=${id}` })
  }
})
