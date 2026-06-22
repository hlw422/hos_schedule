Page({
  data: {
    doctor: {
      name: '医生姓名',
      title: '主治医师'
    }
  },
  goSchedule() {
    wx.switchTab({ url: '/pages/schedule/index' })
  },
  goLeave() {
    wx.navigateTo({ url: '/pages/leave/index' })
  }
})
