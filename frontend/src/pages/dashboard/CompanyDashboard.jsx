import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { listingsApi } from '@/api/listings'
import { applicationsApi } from '@/api/applications'
import StatusBadge from '@/components/StatusBadge'
import SkillTag from '@/components/SkillTag'

export default function CompanyDashboard() {
    const qc = useQueryClient()
    const [selectedListingId, setSelectedListingId] = useState(null)
    const [creatingListing, setCreatingListing] = useState(false)

    const { data: listings, isLoading } = useQuery({
        queryKey: ['my-listings'],
        queryFn: listingsApi.getMyListings,
    })

    const { data: applications } = useQuery({
        queryKey: ['listing-applications', selectedListingId],
        queryFn: () => applicationsApi.getListingApplications(selectedListingId),
        enabled: !!selectedListingId,
    })

    const publishMutation = useMutation({
        mutationFn: listingsApi.publishListing,
        onSuccess: () => qc.invalidateQueries({ queryKey: ['my-listings'] }),
    })

    const closeMutation = useMutation({
        mutationFn: listingsApi.closeListing,
        onSuccess: () => qc.invalidateQueries({ queryKey: ['my-listings'] }),
    })

    const statusMutation = useMutation({
        mutationFn: ({ id, status, comment }) =>
            applicationsApi.changeStatus(id, status, comment),
        onSuccess: () => qc.invalidateQueries({ queryKey: ['listing-applications', selectedListingId] }),
    })

    const statusLabels = {
        draft: 'Черновик',
        active: 'Активна',
        closed: 'Закрыта',
    }

    const statusColors = {
        draft: 'bg-gray-100 text-gray-600',
        active: 'bg-green-100 text-green-700',
        closed: 'bg-red-100 text-red-600',
    }

    return (
        <div className="max-w-6xl mx-auto">
            <div className="flex items-center justify-between mb-6">
                <h1 className="text-2xl font-bold text-gray-900">Кабинет компании</h1>
                <button onClick={() => setCreatingListing(true)} className="btn-primary text-sm px-5 py-2.5">
                    + Создать вакансию
                </button>
            </div>

            {creatingListing && (
                <CreateListingForm
                    onSave={() => {
                        qc.invalidateQueries({ queryKey: ['my-listings'] })
                        setCreatingListing(false)
                    }}
                    onCancel={() => setCreatingListing(false)}
                />
            )}

            <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
                <div className="lg:col-span-2 space-y-3">
                    <h2 className="font-semibold text-gray-700 text-sm uppercase tracking-wide">Вакансии</h2>
                    {isLoading ? (
                        <div className="animate-pulse space-y-3">
                            {[1, 2, 3].map((i) => (
                                <div key={i} className="h-20 bg-white rounded-2xl" />
                            ))}
                        </div>
                    ) : listings?.length === 0 ? (
                        <p className="text-sm text-gray-400 py-4">Нет вакансий. Создайте первую!</p>
                    ) : (
                        listings?.map((listing) => (
                            <div
                                key={listing.id}
                                onClick={() => setSelectedListingId(listing.id)}
                                className={`bg-white rounded-2xl border p-4 cursor-pointer transition-all ${
                                    selectedListingId === listing.id
                                        ? 'border-primary-400 shadow-sm'
                                        : 'border-primary-100 hover:border-primary-300'
                                }`}
                            >
                                <div className="flex items-start justify-between mb-1">
                                    <p className="font-medium text-gray-900 text-sm leading-tight">{listing.title}</p>
                                    <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${statusColors[listing.status]}`}>
                    {statusLabels[listing.status]}
                  </span>
                                </div>
                                <p className="text-xs text-gray-400 mb-2">
                                    {listing.deadline && `до ${new Date(listing.deadline).toLocaleDateString('ru')}`}
                                </p>
                                <div className="flex gap-2">
                                    {listing.status === 'draft' && (
                                        <button
                                            onClick={(e) => {
                                                e.stopPropagation()
                                                publishMutation.mutate(listing.id)
                                            }}
                                            className="text-xs text-primary-600 font-medium hover:underline"
                                        >
                                            Опубликовать
                                        </button>
                                    )}
                                    {listing.status === 'active' && (
                                        <button
                                            onClick={(e) => {
                                                e.stopPropagation()
                                                closeMutation.mutate(listing.id)
                                            }}
                                            className="text-xs text-red-400 font-medium hover:underline"
                                        >
                                            Закрыть
                                        </button>
                                    )}
                                </div>
                            </div>
                        ))
                    )}
                </div>

                <div className="lg:col-span-3">
                    {!selectedListingId ? (
                        <div className="bg-white rounded-2xl border border-primary-100 p-12 text-center h-full flex items-center justify-center">
                            <p className="text-gray-400 text-sm">Выберите вакансию чтобы увидеть отклики</p>
                        </div>
                    ) : (
                        <div className="space-y-3">
                            <h2 className="font-semibold text-gray-700 text-sm uppercase tracking-wide">Отклики</h2>
                            {!applications?.length ? (
                                <div className="bg-white rounded-2xl border border-primary-100 p-10 text-center">
                                    <p className="text-gray-400 text-sm">Нет откликов</p>
                                </div>
                            ) : (
                                applications.map((app) => (
                                    <div key={app.id} className="bg-white rounded-2xl border border-primary-100 p-5">
                                        <div className="flex items-start justify-between mb-3">
                                            <div>
                                                <p className="font-semibold text-gray-900">
                                                    {app.student?.firstName} {app.student?.lastName}
                                                </p>
                                                <p className="text-xs text-gray-400">{app.student?.city}</p>
                                            </div>
                                            <StatusBadge status={app.status} />
                                        </div>

                                        <div className="flex flex-wrap gap-1.5">
                                            {app.student?.skills?.slice(0, 4).map((s) => (
                                                <SkillTag key={s.id} label={s.skill} />
                                            ))}
                                        </div>

                                        {app.coverLetter && (
                                            <p className="text-sm text-gray-600 mt-3 bg-gray-50 rounded-xl px-3 py-2">
                                                {app.coverLetter}
                                            </p>
                                        )}

                                        {app.student?.resumeUrl && (
                                            <a
                                                href={app.student.resumeUrl}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className="inline-block mt-2 text-xs text-primary-600 hover:underline"
                                            >
                                                📎 Скачать резюме
                                            </a>
                                        )}

                                        {app.status !== 'accepted' && app.status !== 'rejected' && (
                                            <div className="flex gap-2 mt-3 flex-wrap">
                                                {app.status === 'applied' && (
                                                    <button
                                                        onClick={() => statusMutation.mutate({ id: app.id, status: 'reviewing' })}
                                                        className="btn-secondary text-xs px-3 py-1.5"
                                                    >
                                                        На рассмотрение
                                                    </button>
                                                )}
                                                {app.status === 'reviewing' && (
                                                    <button
                                                        onClick={() => statusMutation.mutate({ id: app.id, status: 'interview' })}
                                                        className="btn-secondary text-xs px-3 py-1.5"
                                                    >
                                                        Пригласить на интервью
                                                    </button>
                                                )}
                                                <button
                                                    onClick={() => statusMutation.mutate({ id: app.id, status: 'accepted' })}
                                                    className="bg-green-500 text-white text-xs px-3 py-1.5 rounded-lg hover:bg-green-600 transition"
                                                >
                                                    Принять
                                                </button>
                                                <button
                                                    onClick={() => statusMutation.mutate({ id: app.id, status: 'rejected' })}
                                                    className="bg-red-100 text-red-600 text-xs px-3 py-1.5 rounded-lg hover:bg-red-200 transition"
                                                >
                                                    Отказать
                                                </button>
                                            </div>
                                        )}
                                    </div>
                                ))
                            )}
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}

function CreateListingForm({ onSave, onCancel }) {
    const { register, handleSubmit } = useForm()
    const mutation = useMutation({
        mutationFn: listingsApi.createListing,
        onSuccess: onSave,
    })

    return (
        <form
            onSubmit={handleSubmit((v) => mutation.mutate(v))}
            className="bg-white rounded-2xl border border-primary-200 p-6 mb-6 space-y-4"
        >
            <h3 className="font-semibold text-gray-900">Новая вакансия</h3>
            <div className="grid grid-cols-2 gap-4">
                <input {...register('title', { required: true })} placeholder="Название должности *" className="input-base col-span-2" />
                <select {...register('format')} className="input-base">
                    <option value="office">Офис</option>
                    <option value="remote">Удалённо</option>
                    <option value="hybrid">Гибрид</option>
                </select>
                <select {...register('employmentType')} className="input-base">
                    <option value="full_time">Полная занятость</option>
                    <option value="part_time">Частичная</option>
                    <option value="project">Проектная</option>
                </select>
                <input {...register('city')} placeholder="Город" className="input-base" />
                <input {...register('deadline')} type="date" className="input-base" />
                <input {...register('salaryFrom', { valueAsNumber: true })} type="number" placeholder="Зарплата от" className="input-base" />
                <input {...register('salaryTo', { valueAsNumber: true })} type="number" placeholder="Зарплата до" className="input-base" />
            </div>
            <textarea {...register('description', { required: true })} rows={3} placeholder="Описание вакансии *" className="input-base w-full resize-none" />
            <textarea {...register('requirements')} rows={3} placeholder="Требования к кандидату" className="input-base w-full resize-none" />
            <textarea {...register('whatWeOffer')} rows={2} placeholder="Что мы предлагаем" className="input-base w-full resize-none" />
            <div className="flex gap-2">
                <button type="submit" disabled={mutation.isPending} className="btn-primary px-6 py-2.5">
                    {mutation.isPending ? 'Создаём...' : 'Создать как черновик'}
                </button>
                <button type="button" onClick={onCancel} className="btn-secondary px-6 py-2.5">Отмена</button>
            </div>
        </form>
    )
}