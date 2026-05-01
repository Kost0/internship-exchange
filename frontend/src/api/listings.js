import { apiClient } from './client'

export const listingsApi = {
    getListings: (filters) =>
        apiClient.get('/listings', { params: filters }).then((r) => r.data),

    getListing: (id) =>
        apiClient.get(`/listings/${id}`).then((r) => r.data),

    getMyListings: () =>
        apiClient.get('/listings/my').then((r) => r.data),

    createListing: (data) =>
        apiClient.post('/listings', data).then((r) => r.data),

    updateListing: (id, data) =>
        apiClient.put(`/listings/${id}`, data).then((r) => r.data),

    publishListing: (id) =>
        apiClient.post(`/listings/${id}/publish`).then((r) => r.data),

    closeListing: (id) =>
        apiClient.post(`/listings/${id}/close`).then((r) => r.data),

    deleteListing: (id) =>
        apiClient.delete(`/listings/${id}`),
}