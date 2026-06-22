import api from './index'

export const getDepartments = (params) => api.get('/departments', { params })
export const getDepartment = (id) => api.get(`/departments/${id}`)
export const createDepartment = (data) => api.post('/admin/departments', data)
export const updateDepartment = (id, data) => api.put(`/admin/departments/${id}`, data)
export const deleteDepartment = (id) => api.delete(`/admin/departments/${id}`)
