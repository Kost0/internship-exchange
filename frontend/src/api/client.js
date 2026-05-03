import axios from 'axios'
import { useAuthStore } from '../store/auth'

export const apiClient = axios.create({
    baseURL: '/api/v1',
    headers: { 'Content-Type': 'application/json' },
})

apiClient.interceptors.request.use((config) => {
    const token = useAuthStore.getState().accessToken
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

apiClient.interceptors.response.use(
    (res) => res,
    async (error) => {
        const original = error.config

        if (error.response?.status === 401 && !original._retry) {
            original._retry = true

            try {
                const refreshToken = useAuthStore.getState().refreshToken
                const { data } = await axios.post('/api/v1/auth/refresh', { refreshToken })

                useAuthStore.getState().setAccessToken(data.accessToken)
                original.headers.Authorization = `Bearer ${data.accessToken}`

                return apiClient(original)
            } catch (refreshError) {
                useAuthStore.getState().clearAuth()
                window.location.href = '/login'
            }
        }

        return Promise.reject(error)
    }
)