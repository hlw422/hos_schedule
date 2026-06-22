App({
  globalData: {
    userInfo: null,
    token: null,
    baseUrl: 'http://localhost:8080/api/v1'
  },

  onLaunch() {
    this.checkLogin()
  },

  checkLogin() {
    const token = wx.getStorageSync('token')
    if (token) {
      this.globalData.token = token
    }
  },

  login() {
    return new Promise((resolve, reject) => {
      wx.login({
        success: (res) => {
          if (res.code) {
            wx.request({
              url: `${this.globalData.baseUrl}/auth/login`,
              method: 'POST',
              data: { code: res.code },
              success: (resp) => {
                if (resp.data.code === 0) {
                  const { token, user } = resp.data.data
                  this.globalData.token = token
                  this.globalData.userInfo = user
                  wx.setStorageSync('token', token)
                  resolve(user)
                } else {
                  reject(resp.data.message)
                }
              },
              fail: reject
            })
          } else {
            reject(res.errMsg)
          }
        },
        fail: reject
      })
    })
  }
})
