import api from './index'

export const getAppointments = (params) => api.get('/admin/appointments', { params })
export const getAppointmentStats = (params) => api.get('/admin/appointments/stats', { params })
