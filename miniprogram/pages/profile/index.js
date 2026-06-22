const app = getApp()

Page({
  data: {
    userInfo: null
  },

  onShow() {
    this.setData({ userInfo: app.globalData.userInfo })
    if (!this.data.userInfo) {
      this.login()
    }
  },

  async login() {
    try {
      const user = await app.login()
      this.setData({ userInfo: user })
    } catch (e) {
      console.error('Login failed:', e)
    }
  },

  goAppointments() {
    wx.switchTab({ url: '/pages/appointment/list' })
  },

  goPatients() {
    wx.navigateTo({ url: '/pages/patients/list' })
  },

  goAbout() {
    wx.showModal({
      title: '关于我们',
      content: '医院预约挂号小程序 v1.0.0',
      showCancel: false
    })
  },

  callService() {
    wx.makePhoneCall({
      phoneNumber: '400-123-4567',
      fail: () => {}
    })
  },

  logout() {
    wx.showModal({
      title: '退出登录',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          app.globalData.token = null
          app.globalData.userInfo = null
          wx.removeStorageSync('token')
          this.setData({ userInfo: null })
          wx.showToast({ title: '已退出', icon: 'success' })
        }
      }
    })
  }
})
