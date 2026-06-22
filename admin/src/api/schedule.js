import api from './index'

export const getSchedules = (params) => api.get('/schedules', { params })
export const getSchedule = (id) => api.get(`/schedules/${id}`)
export const createSchedule = (data) => api.post('/admin/schedules', data)
export const batchCreateSchedule = (data) => api.post('/admin/schedules/batch', data)
export const updateSchedule = (id, data) => api.put(`/admin/schedules/${id}`, data)
export const deleteSchedule = (id) => api.delete(`/admin/schedules/${id}`)
