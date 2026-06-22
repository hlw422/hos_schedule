const api = require('../../utils/api')

Page({
  data: {
    doctors: [],
    selectedDoctor: {},
    date: '',
    period: '',
    periods: ['上午', '下午', '全天'],
    reason: '',
    today: '',
    canSubmit: false
  },

  onLoad() {
    this.setData({
      today: new Date().toISOString().split('T')[0]
    })
    this.loadDoctors()
  },

  async loadDoctors() {
    try {
      const doctors = await api.getDoctors()
      this.setData({ doctors: doctors || [] })
    } catch (err) {
      console.error('加载医生列表失败:', err)
    }
  },

  onDoctorChange(e) {
    const idx = e.detail.value
    this.setData({
      selectedDoctor: this.data.doctors[idx]
    })
    this.checkCanSubmit()
  },

  onDateChange(e) {
    this.setData({ date: e.detail.value })
    this.checkCanSubmit()
  },

  onPeriodChange(e) {
    const idx = e.detail.value
    this.setData({ period: this.data.periods[idx] })
    this.checkCanSubmit()
  },

  onReasonInput(e) {
    this.setData({ reason: e.detail.value })
  },

  checkCanSubmit() {
    const { selectedDoctor, date, period } = this.data
    this.setData({
      canSubmit: !!(selectedDoctor.id && date && period)
    })
  },

  async submit() {
    const { selectedDoctor, date, period, reason } = this.data
    const periodMap = { '上午': 'morning', '下午': 'afternoon', '全天': 'all' }

    wx.showModal({
      title: '确认停诊',
      content: `${selectedDoctor.name} ${date} ${period}`,
      success: async (res) => {
        if (res.confirm) {
          try {
            await api.stopSchedule({
              doctor_id: selectedDoctor.id,
              date: date,
              period: periodMap[period],
              reason: reason
            })
            wx.showToast({ title: '停诊设置成功', icon: 'success' })
            this.setData({ date: '', period: '', reason: '', canSubmit: false })
          } catch (err) {
            console.error('停诊设置失败:', err)
          }
        }
      }
    })
  }
})
