import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { applicationsApi } from '@/api/applications'
import StatusBadge from '@/components/StatusBadge'

export default function StudentDashboard() {
    const qc = useQueryClient()

    const { data: applications, isLoading } = useQuery({
        queryKey: ['my-applications'],
        queryFn: applicationsApi.getMyApplications,
    })

    const withdraw = useMutation({
        mutationFn: applicationsApi.withdrawApplication,
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['my-applications'] })
        },
    })

    if (isLoading) {
        return (
            <div className="animate-pulse space-y-4">
                {Array.from({ length: 4 }).map((_, i) => (
                    <div key={i} className="h-24 bg-white rounded-2xl" />
                ))}
            </div>
        )
    }

    return (
        <div className="max-w-3xl mx-auto">
            <div className="flex items-center justify-between mb-6">
                <h1 className="text-2xl font-bold text-gray-900">Мои отклики</h1>
                <Link to="/" className="btn-primary text-sm px-5 py-2.5">
                    Найти вакансии
                </Link>
            </div>

            {!applications?.length ? (
                <div className="bg-white rounded-2xl border border-primary-100 p-12 text-center">
                    <p className="text-gray-400 mb-4">Вы ещё не откликались на вакансии</p>
                    <Link to="/" className="btn-primary text-sm px-6 py-2.5">
                        Перейти к вакансиям
                    </Link>
                </div>
            ) : (
                <div className="space-y-3">
                    {applications.map((app) => (
                        <div key={app.id} className="bg-white rounded-2xl border border-primary-100 p-5">
                            <div className="flex items-start justify-between">
                                <div className="flex-1">
                                    <Link
                                        to={`/listings/${app.listingId}`}
                                        className="font-semibold text-gray-900 hover:text-primary-600 transition-colors"
                                    >
                                        {app.listing?.title ?? 'Вакансия'}
                                    </Link>
                                    <p className="text-sm text-gray-500 mt-0.5">{app.listing?.company?.name}</p>
                                    <p className="text-xs text-gray-400 mt-1">
                                        Отклик отправлен {new Date(app.createdAt).toLocaleDateString('ru')}
                                    </p>
                                </div>
                                <div className="flex flex-col items-end gap-2">
                                    <StatusBadge status={app.status} />

                                    {app.status === 'applied' && (
                                        <button
                                            onClick={() => {
                                                if (confirm('Вы уверены, что хотите отозвать отклик?')) {
                                                    withdraw.mutate(app.id)
                                                }
                                            }}
                                            disabled={withdraw.isPending}
                                            className="text-xs text-red-400 hover:text-red-600 transition disabled:opacity-50"
                                        >
                                            {withdraw.isPending ? 'Загрузка...' : 'Отозвать'}
                                        </button>
                                    )}
                                </div>
                            </div>

                            {app.events?.length > 0 && app.events[app.events.length - 1].comment && (
                                <div className="mt-3 bg-primary-50 rounded-xl px-4 py-2.5">
                                    <p className="text-xs text-gray-500 mb-0.5">Комментарий компании:</p>
                                    <p className="text-sm text-gray-700">
                                        {app.events[app.events.length - 1].comment}
                                    </p>
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            )}
        </div>
    )
}