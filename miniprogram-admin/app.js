App({
  globalData: {
    token: null,
    baseUrl: 'http://localhost:8080/api/v1'
  },
  onLaunch() {
    const token = wx.getStorageSync('admin_token')
    if (token) this.globalData.token = token
  }
})
