import { apiClient } from './client'

export const profileApi = {
    getMyStudentProfile: () =>
        apiClient.get('/profile/student').then((r) => r.data),

    getStudentProfile: (id) =>
        apiClient.get(`/profile/student/${id}`).then((r) => r.data),

    updateStudentProfile: (data) =>
        apiClient.put('/profile/student', data).then((r) => r.data),

    uploadAvatar: (file) => {
        const form = new FormData()
        form.append('avatar', file)
        return apiClient
            .post('/profile/student/avatar', form, {
                headers: { 'Content-Type': 'multipart/form-data' },
            })
            .then((r) => r.data)
    },

    uploadResume: (file) => {
        const form = new FormData()
        form.append('resume', file)
        return apiClient
            .post('/profile/student/resume', form, {
                headers: { 'Content-Type': 'multipart/form-data' },
            })
            .then((r) => r.data)
    },

    getResumeUrl: (id) =>
        apiClient
            .get(`/profile/student/${id}/resume`)
            .then((r) => r.data),

    addEducation: (data) =>
        apiClient.post('/profile/student/education', data).then((r) => r.data),

    updateEducation: (id, data) =>
        apiClient.put(`/profile/student/education/${id}`, data).then((r) => r.data),

    deleteEducation: (id) =>
        apiClient.delete(`/profile/student/education/${id}`),

    addExperience: (data) =>
        apiClient.post('/profile/student/experience', data).then((r) => r.data),

    updateExperience: (id, data) =>
        apiClient.put(`/profile/student/experience/${id}`, data).then((r) => r.data),

    deleteExperience: (id) =>
        apiClient.delete(`/profile/student/experience/${id}`),

    addProject: (data) =>
        apiClient.post('/profile/student/projects', data).then((r) => r.data),

    updateProject: (id, data) =>
        apiClient.put(`/profile/student/projects/${id}`, data).then((r) => r.data),

    deleteProject: (id) =>
        apiClient.delete(`/profile/student/projects/${id}`),

    getMyCompanyProfile: () =>
        apiClient.get('/profile/company').then((r) => r.data),

    getCompanyProfile: (id) =>
        apiClient.get(`/profile/company/${id}`).then((r) => r.data),

    updateCompanyProfile: (data) =>
        apiClient.put('/profile/company', data).then((r) => r.data),

    uploadLogo: (file) => {
        const form = new FormData()
        form.append('logo', file)
        return apiClient
            .post('/profile/company/logo', form, {
                headers: { 'Content-Type': 'multipart/form-data' },
            })
            .then((r) => r.data)
    },
}