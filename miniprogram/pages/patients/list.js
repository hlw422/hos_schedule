const api = require('../../utils/api')

Page({
  data: {
    patients: [],
    loading: false
  },

  onShow() {
    this.loadPatients()
  },

  async loadPatients() {
    this.setData({ loading: true })
    try {
      const patients = await api.getPatients()
      this.setData({ patients })
    } catch (e) {
      console.error('Failed to load patients:', e)
    } finally {
      this.setData({ loading: false })
    }
  },

  goAdd() {
    wx.navigateTo({ url: '/pages/patients/edit' })
  },

  goEdit(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({ url: `/pages/patients/edit?id=${id}` })
  },

  async setDefault(e) {
    const id = e.currentTarget.dataset.id
    try {
      await api.setDefaultPatient(id)
      wx.showToast({ title: '设置成功', icon: 'success' })
      this.loadPatients()
    } catch (e) {
      wx.showToast({ title: e || '设置失败', icon: 'none' })
    }
  },

  deletePatient(e) {
    const id = e.currentTarget.dataset.id
    wx.showModal({
      title: '删除就诊人',
      content: '确定要删除该就诊人吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            await api.deletePatient(id)
            wx.showToast({ title: '删除成功', icon: 'success' })
            this.loadPatients()
          } catch (e) {
            wx.showToast({ title: e || '删除失败', icon: 'none' })
          }
        }
      }
    })
  }
})
