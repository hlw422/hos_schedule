const api = require('../../utils/api')

Page({
  data: {
    id: '',
    hospital: {},
    campuses: [],
    currentCampus: 0,
    departments: []
  },

  onLoad(options) {
    this.setData({ id: options.id })
    this.loadHospital()
  },

  async loadHospital() {
    try {
      const hospital = await api.getHospital(this.data.id)
      this.setData({ hospital })
      this.loadCampuses()
    } catch (e) {
      console.error('Failed to load hospital:', e)
    }
  },

  async loadCampuses() {
    try {
      const campuses = await api.getCampuses(this.data.id)
      this.setData({ campuses })
      if (campuses.length > 0) {
        this.loadDepartments(campuses[0].id)
      }
    } catch (e) {
      console.error('Failed to load campuses:', e)
    }
  },

  async loadDepartments(campusId) {
    try {
      const departments = await api.getDepartments({ hospital_id: this.data.id, campus_id: campusId })
      this.setData({ departments })
    } catch (e) {
      console.error('Failed to load departments:', e)
    }
  },

  switchCampus(e) {
    const index = e.currentTarget.dataset.index
    this.setData({ currentCampus: index })
    this.loadDepartments(this.data.campuses[index].id)
  },

  goDeptDoctors(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({ url: `/pages/doctors/list?department_id=${id}` })
  }
})
