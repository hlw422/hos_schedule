const api = require('../../utils/api')

Page({
  data: {
    weekDays: ['日', '一', '二', '三', '四', '五', '六'],
    dates: [],
    schedules: [],
    selectedDate: ''
  },
  onLoad() {
    this.initCalendar()
  },
  initCalendar() {
    const now = new Date()
    const year = now.getFullYear()
    const month = now.getMonth()
    const today = now.getDate()
    const firstDay = new Date(year, month, 1).getDay()
    const daysInMonth = new Date(year, month + 1, 0).getDate()

    const dates = []
    for (let i = 0; i < firstDay; i++) {
      dates.push({ day: '', date: '', empty: true })
    }
    for (let d = 1; d <= daysInMonth; d++) {
      const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
      dates.push({
        day: d,
        date: dateStr,
        isToday: d === today,
        hasSchedule: false
      })
    }
    this.setData({ dates })
    this.loadSchedules(year, month + 1)
  },
  async loadSchedules(year, month) {
    try {
      const data = await api.getSchedules({ year, month })
      if (data) {
        const dates = this.data.dates.map(d => {
          if (d.date && data.find(s => s.date === d.date)) {
            d.hasSchedule = true
          }
          return d
        })
        this.setData({ dates })
      }
    } catch (e) {
      console.error(e)
    }
  },
  selectDate(e) {
    const date = e.currentTarget.dataset.date
    if (!date) return
    this.setData({ selectedDate: date })
    // Could load detailed schedule for this date
  }
})
