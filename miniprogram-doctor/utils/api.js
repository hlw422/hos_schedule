const app = getApp()

const request = (options) => new Promise((resolve, reject) => {
  wx.request({
    url: `${app.globalData.baseUrl}${options.url}`,
    method: options.method || 'GET',
    data: options.data,
    header: { 'Authorization': `Bearer ${app.globalData.token}` },
    success: (res) => res.data.code === 0 ? resolve(res.data.data) : reject(res.data.message),
    fail: reject
  })
})

module.exports = {
  getSchedules: (params) => request({ url: '/doctor/schedules', data: params }),
  getTodayAppointments: () => request({ url: '/doctor/appointments' }),
  requestLeave: (data) => request({ url: '/doctor/leave', method: 'POST', data })
}
