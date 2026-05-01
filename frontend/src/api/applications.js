import { apiClient } from './client'

export const applicationsApi = {
    apply: (listingId, coverLetter) =>
        apiClient
            .post('/applications', { listingId, coverLetter })
            .then((r) => r.data),

    getMyApplications: () =>
        apiClient.get('/applications/my').then((r) => r.data),

    withdrawApplication: (id) =>
        apiClient.delete(`/applications/${id}`),

    getListingApplications: (listingId) =>
        apiClient
            .get(`/listings/${listingId}/applications`)
            .then((r) => r.data),

    changeStatus: (id, status, comment) =>
        apiClient
            .put(`/applications/${id}/status`, { status, comment })
            .then((r) => r.data),

    getApplicationHistory: (id) =>
        apiClient
            .get(`/applications/${id}/history`)
            .then((r) => r.data),
}