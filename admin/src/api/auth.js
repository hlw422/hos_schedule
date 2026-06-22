import api from './index'

export const login = (data) => api.post('/auth/login', data)
export const getProfile = () => api.get('/me')
