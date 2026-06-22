import api from './index'

export const getHospitals = () => api.get('/hospitals')
export const getHospital = (id) => api.get(`/hospitals/${id}`)
export const updateHospital = (id, data) => api.put(`/admin/hospitals/${id}`, data)
export const getCampuses = (id) => api.get(`/hospitals/${id}/campuses`)
export const addCampus = (data) => api.post('/admin/hospitals/campuses', data)
