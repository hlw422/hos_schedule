const app = getApp()

const request = (options) => new Promise((resolve, reject) => {
  wx.request({
    url: `${app.globalData.baseUrl}${options.url}`,
    method: options.method || 'GET',
    data: options.data,
    header: {
      'Authorization': `Bearer ${app.globalData.token}`,
      'Content-Type': 'application/json'
    },
    success: (res) => {
      if (res.data.code === 0) {
        resolve(res.data.data)
      } else {
        wx.showToast({ title: res.data.message || '请求失败', icon: 'none' })
        reject(res.data.message)
      }
    },
    fail: (err) => {
      wx.showToast({ title: '网络错误', icon: 'none' })
      reject(err)
    }
  })
})

module.exports = {
  getAppointmentStats: (params) => request({ url: '/admin/appointments/stats', data: params }),
  getAppointments: (params) => request({ url: '/admin/appointments', data: params }),
  stopSchedule: (data) => request({ url: '/admin/schedules/stop', method: 'POST', data }),
  getDoctors: () => request({ url: '/admin/doctors' }),
  getStatsTrend: (params) => request({ url: '/admin/appointments/trend', data: params })
}
