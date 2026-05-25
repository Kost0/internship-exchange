import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { listingsApi } from '../../api/listings'
import JobCard from '../../components/JobCard'

const formats = [
    { value: '', label: 'Все' },
    { value: 'remote', label: 'Удалённо' },
    { value: 'office', label: 'Офис' },
    { value: 'hybrid', label: 'Гибрид' },
]

export default function ListingsPage() {
    const [filters, setFilters] = useState({ page: 1, limit: 12 })
    const [query, setQuery] = useState('')

    const { data, isLoading } = useQuery({
        queryKey: ['listings', filters],
        queryFn: () => listingsApi.getListings(filters),
    })

    const handleSearch = (e) => {
        e.preventDefault()
        setFilters(f => ({ ...f, query, page: 1 }))
    }

    return (
        <div>
            <h1 style={{ marginBottom: 16 }}>Вакансии для стажировок</h1>

            <form onSubmit={handleSearch} style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
                <input
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder="Поиск..."
                    style={{ flex: 1, padding: '8px 12px', border: '1px solid #ccc', fontSize: 14 }}
                />
                <button type="submit" className="btn-primary">Найти</button>
            </form>

            <div style={{ display: 'flex', gap: 8, marginBottom: 16, flexWrap: 'wrap' }}>
                <span style={{ fontSize: 13, color: '#555', marginRight: 4, alignSelf: 'center' }}>Формат:</span>
                {formats.map(f => (
                    <button
                        key={f.value}
                        onClick={() => setFilters(prev => ({ ...prev, format: f.value || undefined, page: 1 }))}
                        style={{
                            padding: '4px 12px', fontSize: 13, cursor: 'pointer',
                            border: '1px solid #ccc',
                            background: (filters.format === f.value || (!filters.format && f.value === '')) ? '#3e85dc' : 'white',
                            color: (filters.format === f.value || (!filters.format && f.value === '')) ? 'white' : '#333',
                        }}
                    >
                        {f.label}
                    </button>
                ))}
            </div>

            <div style={{ fontSize: 13, color: '#555', marginBottom: 12 }}>
                Найдено: {data?.total ?? 0}
            </div>

            {isLoading ? (
                <div>Загрузка...</div>
            ) : (
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: 16 }}>
                    {data?.items?.map(listing => (
                        <JobCard key={listing.id} listing={listing} />
                    ))}
                </div>
            )}

            {(data?.total ?? 0) > (data?.limit ?? 12) && (
                <div style={{ display: 'flex', gap: 8, marginTop: 24, justifyContent: 'center', alignItems: 'center' }}>
                    <button
                        disabled={filters.page === 1}
                        onClick={() => setFilters(f => ({ ...f, page: f.page - 1 }))}
                        className="btn-secondary"
                        style={{ opacity: filters.page === 1 ? 0.4 : 1 }}
                    >
                        ← Назад
                    </button>
                    <span style={{ fontSize: 13 }}>
                        {filters.page} / {Math.ceil((data?.total ?? 0) / (data?.limit ?? 12))}
                    </span>
                    <button
                        disabled={filters.page >= Math.ceil((data?.total ?? 0) / (data?.limit ?? 12))}
                        onClick={() => setFilters(f => ({ ...f, page: f.page + 1 }))}
                        className="btn-secondary"
                        style={{ opacity: filters.page >= Math.ceil((data?.total ?? 0) / (data?.limit ?? 12)) ? 0.4 : 1 }}
                    >
                        Вперёд →
                    </button>
                </div>
            )}
        </div>
    )
}