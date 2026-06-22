const api = require('../../utils/api')

Page({
  data: {
    banners: [
      { id: 1, image: '/images/banner1.png' },
      { id: 2, image: '/images/banner2.png' }
    ],
    hotDepts: [
      { id: 1, name: '内科', icon: '/images/dept-internal.png' },
      { id: 2, name: '外科', icon: '/images/dept-surgery.png' },
      { id: 3, name: '儿科', icon: '/images/dept-pediatrics.png' },
      { id: 4, name: '妇科', icon: '/images/dept-gynecology.png' },
      { id: 5, name: '骨科', icon: '/images/dept-orthopedics.png' },
      { id: 6, name: '眼科', icon: '/images/dept-ophthalmology.png' },
      { id: 7, name: '口腔科', icon: '/images/dept-dental.png' },
      { id: 8, name: '皮肤科', icon: '/images/dept-dermatology.png' }
    ],
    doctors: []
  },

  onLoad() {
    this.loadDoctors()
  },

  async loadDoctors() {
    try {
      const doctors = await api.getDoctors({ limit: 5 })
      this.setData({ doctors })
    } catch (e) {
      console.error('Failed to load doctors:', e)
    }
  },

  goSearch() {
    wx.navigateTo({ url: '/pages/hospitals/list' })
  },

  goHospitals() {
    wx.navigateTo({ url: '/pages/hospitals/list' })
  },

  goAppointments() {
    wx.switchTab({ url: '/pages/appointment/list' })
  },

  goPatients() {
    wx.navigateTo({ url: '/pages/patients/list' })
  },

  goDeptDoctors(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({ url: `/pages/doctors/list?department_id=${id}` })
  },

  goDoctor(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({ url: `/pages/doctors/detail?id=${id}` })
  }
})
