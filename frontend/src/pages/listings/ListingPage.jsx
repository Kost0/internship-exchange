import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { listingsApi } from '../../api/listings'
import JobCard from '../../components/JobCard'

const formats = [
    { value: '', label: 'Все форматы' },
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
        setFilters((f) => ({ ...f, query, page: 1 }))
    }

    return (
        <div>
            <div className="bg-primary-500 rounded-3xl px-10 py-12 mb-8 text-white">
                <h1 className="text-3xl font-bold mb-2">Найди стажировку своей мечты</h1>
                <p className="text-primary-100 mb-6">Тысячи компаний ищут молодых специалистов прямо сейчас</p>
                <form onSubmit={handleSearch} className="flex gap-3 max-w-2xl">
                    <input
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                        placeholder="Поиск по вакансиям, компаниям, навыкам..."
                        className="flex-1 px-5 py-3 rounded-xl text-gray-900 text-sm border border-white/40 focus:outline-none focus:ring-2 focus:ring-white"
                    />
                    <button
                        type="submit"
                        className="bg-white text-primary-600 font-semibold px-6 py-3 rounded-xl hover:bg-primary-50 transition-colors"
                    >
                        Найти
                    </button>
                </form>
            </div>

            <div className="flex items-center gap-3 mb-6 flex-wrap">
                <span className="text-sm font-medium text-gray-500">Формат:</span>
                {formats.map((f) => (
                    <button
                        key={f.value}
                        onClick={() => setFilters((prev) => ({ ...prev, format: f.value || undefined, page: 1 }))}
                        className={`text-sm px-4 py-1.5 rounded-full font-medium transition-colors ${
                            filters.format === f.value || (!filters.format && f.value === '')
                                ? 'bg-primary-500 text-white'
                                : 'bg-white text-gray-600 border border-gray-200 hover:border-primary-300'
                        }`}
                    >
                        {f.label}
                    </button>
                ))}
            </div>

            <div className="flex items-center justify-between mb-4">
                <h2 className="font-semibold text-gray-900">
                    {`Найдено вакансий: ${data?.total ?? 0}`}
                </h2>
            </div>

            {isLoading ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                    {Array.from({ length: 6 }).map((_, i) => (
                        <div key={i} className="bg-white rounded-2xl border border-primary-100 p-5 h-48 animate-pulse" />
                    ))}
                </div>
            ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                    {data.items?.map((listing) => (
                        <JobCard key={listing.id} listing={listing} />
                    ))}
                </div>
            )}

            {/* Pagination */}
            {data && data.total > data.limit && (
                <div className="flex justify-center gap-2 mt-8">
                    <button
                        disabled={filters.page === 1}
                        onClick={() => setFilters((f) => ({ ...f, page: (f.page ?? 1) - 1 }))}
                        className="px-4 py-2 rounded-xl border border-gray-200 text-sm font-medium disabled:opacity-40 hover:border-primary-300 transition"
                    >
                        ← Назад
                    </button>
                    <span className="px-4 py-2 text-sm text-gray-500">
            Страница {filters.page} из {Math.ceil(data.total / data.limit)}
          </span>
                    <button
                        disabled={(filters.page ?? 1) >= Math.ceil(data.total / data.limit)}
                        onClick={() => setFilters((f) => ({ ...f, page: (f.page ?? 1) + 1 }))}
                        className="px-4 py-2 rounded-xl border border-gray-200 text-sm font-medium disabled:opacity-40 hover:border-primary-300 transition"
                    >
                        Вперёд →
                    </button>
                </div>
            )}
        </div>
    )
}