const api = require('../../utils/api')

Page({
  data: {
    date: '',
    period: '',
    reason: '',
    periods: ['上午', '下午', '全天']
  },
  onDateChange(e) {
    this.setData({ date: e.detail.value })
  },
  onPeriodChange(e) {
    this.setData({ period: this.data.periods[e.detail.value] })
  },
  onReasonInput(e) {
    this.setData({ reason: e.detail.value })
  },
  async submit() {
    const { date, period, reason } = this.data
    if (!date || !period) {
      wx.showToast({ title: '请选择日期和时段', icon: 'none' })
      return
    }
    try {
      const periodMap = { '上午': 'MORNING', '下午': 'AFTERNOON', '全天': 'FULL_DAY' }
      await api.requestLeave({ date, time_period: periodMap[period], reason })
      wx.showToast({ title: '提交成功', icon: 'success' })
      setTimeout(() => wx.navigateBack(), 1500)
    } catch (e) {
      wx.showToast({ title: '提交失败', icon: 'none' })
    }
  }
})
