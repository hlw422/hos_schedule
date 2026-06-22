const api = require('../../utils/api')

Page({
  data: {
    keyword: '',
    hospitals: [],
    loading: false
  },

  onLoad() {
    this.loadHospitals()
  },

  async loadHospitals() {
    this.setData({ loading: true })
    try {
      const hospitals = await api.getHospitals()
      this.setData({ hospitals })
    } catch (e) {
      console.error('Failed to load hospitals:', e)
    } finally {
      this.setData({ loading: false })
    }
  },

  onInput(e) {
    this.setData({ keyword: e.detail.value })
  },

  onSearch() {
    this.loadHospitals()
  },

  goDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({ url: `/pages/hospitals/detail?id=${id}` })
  }
})
