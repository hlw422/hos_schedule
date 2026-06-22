const app = getApp()

const request = (options) => {
  return new Promise((resolve, reject) => {
    const { url, method = 'GET', data } = options
    const token = app.globalData.token
    
    wx.request({
      url: `${app.globalData.baseUrl}${url}`,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (res) => {
        if (res.data.code === 0) {
          resolve(res.data.data)
        } else if (res.data.code === 401) {
          app.login().then(() => {
            request(options).then(resolve).catch(reject)
          }).catch(reject)
        } else {
          reject(res.data.message)
        }
      },
      fail: reject
    })
  })
}

module.exports = {
  getHospitals: () => request({ url: '/hospitals' }),
  getHospital: (id) => request({ url: `/hospitals/${id}` }),
  getCampuses: (id) => request({ url: `/hospitals/${id}/campuses` }),
  getDepartments: (params) => request({ url: '/departments', data: params }),
  getDepartment: (id) => request({ url: `/departments/${id}` }),
  getDoctors: (params) => request({ url: '/doctors', data: params }),
  getDoctor: (id) => request({ url: `/doctors/${id}` }),
  getSchedules: (params) => request({ url: '/schedules', data: params }),
  createAppointment: (data) => request({ url: '/appointments', method: 'POST', data }),
  getAppointments: (params) => request({ url: '/appointments', data: params }),
  getAppointment: (id) => request({ url: `/appointments/${id}` }),
  cancelAppointment: (id, reason) => request({ url: `/appointments/${id}/cancel`, method: 'PUT', data: { reason } }),
  getPatients: () => request({ url: '/patients' }),
  createPatient: (data) => request({ url: '/patients', method: 'POST', data }),
  updatePatient: (id, data) => request({ url: `/patients/${id}`, method: 'PUT', data }),
  deletePatient: (id) => request({ url: `/patients/${id}`, method: 'DELETE' }),
  setDefaultPatient: (id) => request({ url: `/patients/${id}/default`, method: 'PUT' })
}
