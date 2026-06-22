App({
  globalData: {
    userInfo: null,
    token: null,
    baseUrl: 'http://localhost:8080/api/v1'
  },
  onLaunch() {
    const token = wx.getStorageSync('doctor_token')
    if (token) this.globalData.token = token
  },
  login() {
    return new Promise((resolve, reject) => {
      wx.login({
        success: (res) => {
          wx.request({
            url: `${this.globalData.baseUrl}/auth/login`,
            method: 'POST',
            data: { code: res.code },
            success: (resp) => {
              if (resp.data.code === 0) {
                this.globalData.token = resp.data.data.token
                wx.setStorageSync('doctor_token', resp.data.data.token)
                resolve(resp.data.data)
              } else reject(resp.data.message)
            },
            fail: reject
          })
        }
      })
    })
  }
})
