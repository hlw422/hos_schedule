const api = require('../../utils/api')

Page({
  data: {
    scheduleId: '',
    doctorId: '',
    doctor: {},
    schedule: {},
    patients: [],
    selectedPatient: null,
    symptom: '',
    submitting: false
  },

  onLoad(options) {
    this.setData({
      scheduleId: options.schedule_id,
      doctorId: options.doctor_id
    })
    this.loadData()
  },

  onShow() {
    this.loadPatients()
  },

  async loadData() {
    try {
      const [doctor, schedule] = await Promise.all([
        api.getDoctor(this.data.doctorId),
        api.getSchedules({ doctor_id: this.data.doctorId }).then(schedules => 
          schedules.find(s => s.id == this.data.scheduleId)
        )
      ])
      this.setData({ doctor, schedule })
    } catch (e) {
      console.error('Failed to load data:', e)
    }
  },

  async loadPatients() {
    try {
      const patients = await api.getPatients()
      this.setData({ patients })
      const defaultPatient = patients.find(p => p.is_default)
      if (defaultPatient) {
        this.setData({ selectedPatient: defaultPatient.id })
      }
    } catch (e) {
      console.error('Failed to load patients:', e)
    }
  },

  selectPatient(e) {
    const id = e.currentTarget.dataset.id
    this.setData({ selectedPatient: id })
  },

  onSymptomInput(e) {
    this.setData({ symptom: e.detail.value })
  },

  goAddPatient() {
    wx.navigateTo({ url: '/pages/patients/edit' })
  },

  async submit() {
    if (!this.data.selectedPatient || this.data.submitting) return
    this.setData({ submitting: true })
    try {
      const result = await api.createAppointment({
        schedule_id: this.data.scheduleId,
        patient_id: this.data.selectedPatient,
        symptom: this.data.symptom
      })
      wx.redirectTo({
        url: `/pages/appointment/success?id=${result.id}`
      })
    } catch (e) {
      wx.showToast({ title: e || '预约失败', icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  }
})
