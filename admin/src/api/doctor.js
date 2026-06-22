import api from './index'

export const getDoctors = (params) => api.get('/doctors', { params })
export const getDoctor = (id) => api.get(`/doctors/${id}`)
export const createDoctor = (data) => api.post('/admin/doctors', data)
export const updateDoctor = (id, data) => api.put(`/admin/doctors/${id}`, data)
export const updateDoctorStatus = (id, status) => api.put(`/admin/doctors/${id}/status`, { status })
