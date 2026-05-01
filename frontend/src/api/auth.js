import { apiClient } from './client'

export const authApi = {
    register: (data) =>
        apiClient.post('/auth/register', data).then((r) => r.data),

    login: (data) =>
        apiClient.post('/auth/login', data).then((r) => r.data),

    refresh: (refreshToken) =>
        apiClient
            .post('/auth/refresh', { refreshToken })
            .then((r) => r.data),
}