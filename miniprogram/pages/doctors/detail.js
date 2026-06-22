const api = require('../../utils/api')

Page({
  data: {
    id: '',
    doctor: {},
    dates: [],
    currentDate: 0,
    schedules: [],
    selectedSchedule: null
  },

  onLoad(options) {
    this.setData({ id: options.id })
    this.loadDoctor()
    this.initDates()
  },

  async loadDoctor() {
    try {
      const doctor = await api.getDoctor(this.data.id)
      this.setData({ doctor })
    } catch (e) {
      console.error('Failed to load doctor:', e)
    }
  },

  initDates() {
    const dates = []
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    for (let i = 0; i < 7; i++) {
      const date = new Date()
      date.setDate(date.getDate() + i)
      dates.push({
        date: `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`,
        weekday: i === 0 ? '今天' : i === 1 ? '明天' : weekdays[date.getDay()],
        day: String(date.getDate()).padStart(2, '0')
      })
    }
    this.setData({ dates })
    this.loadSchedules()
  },

  selectDate(e) {
    const index = e.currentTarget.dataset.index
    this.setData({ currentDate: index, selectedSchedule: null })
    this.loadSchedules()
  },

  async loadSchedules() {
    try {
      const date = this.data.dates[this.data.currentDate].date
      const schedules = await api.getSchedules({
        doctor_id: this.data.id,
        date: date
      })
      this.setData({ schedules })
    } catch (e) {
      console.error('Failed to load schedules:', e)
    }
  },

  selectSchedule(e) {
    const { id, remaining } = e.currentTarget.dataset
    if (remaining === 0) return
    this.setData({ selectedSchedule: id })
  },

  goConfirm() {
    if (!this.data.selectedSchedule) return
    wx.navigateTo({
      url: `/pages/appointment/confirm?schedule_id=${this.data.selectedSchedule}&doctor_id=${this.data.id}`
    })
  }
})
