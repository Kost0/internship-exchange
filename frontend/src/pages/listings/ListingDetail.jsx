import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { listingsApi } from '../../api/listings'
import { applicationsApi } from '../../api/applications'
import { useAuthStore } from '../../store/auth'
import SkillTag from '../../components/SkillTag'

const formatLabels = { office: 'Офис', remote: 'Удалённо', hybrid: 'Гибрид' }
const employmentLabels = { full_time: 'Полная занятость', part_time: 'Частичная', project: 'Проектная' }

export default function ListingDetail() {
    const { id } = useParams()
    const { user } = useAuthStore()
    const qc = useQueryClient()
    const [applied, setApplied] = useState(false)

    const { data: listing, isLoading } = useQuery({
        queryKey: ['listing', id],
        queryFn: () => listingsApi.getListing(id),
        enabled: !!id,
    })

    const { register, handleSubmit, formState: { errors } } = useForm()

    const applyMutation = useMutation({
        mutationFn: ({ coverLetter }) =>
            applicationsApi.apply(id, coverLetter),
        onSuccess: () => {
            setApplied(true)
            qc.invalidateQueries({ queryKey: ['my-applications'] })
        },
    })

    if (isLoading) {
        return (
            <div className="animate-pulse space-y-4">
                <div className="h-48 bg-white rounded-3xl" />
                <div className="h-96 bg-white rounded-3xl" />
            </div>
        )
    }

    if (!listing) return <p className="text-center text-gray-500 py-20">Вакансия не найдена</p>

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Main content */}
            <div className="lg:col-span-2 space-y-6">
                {/* Header */}
                <div className="bg-white rounded-2xl border border-primary-100 p-6">
                    <div className="flex items-start gap-4 mb-4">
                        <div className="w-16 h-16 rounded-xl bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-xl flex-shrink-0">
                            {listing.company?.name?.[0] ?? '?'}
                        </div>
                        <div>
                            <h1 className="text-xl font-bold text-gray-900">{listing.title}</h1>
                            <Link
                                to={`/companies/${listing.companyId}`}
                                className="text-sm text-primary-600 font-medium hover:underline"
                            >
                                {listing.company?.name}
                            </Link>
                            <div className="flex items-center gap-2 text-sm text-gray-500 mt-1">
                                {listing.city && <span>{listing.city}</span>}
                                <span>·</span>
                                <span>{formatLabels[listing.format]}</span>
                                <span>·</span>
                                <span>{employmentLabels[listing.employmentType]}</span>
                            </div>
                        </div>
                    </div>

                    {(listing.salaryFrom || listing.salaryTo) && (
                        <p className="text-lg font-bold text-primary-600 mb-3">
                            {listing.salaryFrom && `от ${listing.salaryFrom.toLocaleString('ru')} ₽`}
                            {listing.salaryFrom && listing.salaryTo && ' — '}
                            {listing.salaryTo && `до ${listing.salaryTo.toLocaleString('ru')} ₽`}
                        </p>
                    )}

                    <div className="flex flex-wrap gap-1.5 mb-3">
                        {listing.skills.map((s) => (
                            <SkillTag key={s.id} label={s.skill} variant={s.isRequired ? 'primary' : 'gray'} />
                        ))}
                    </div>

                    {listing.deadline && (
                        <p className="text-xs text-gray-400">Дедлайн: {new Date(listing.deadline).toLocaleDateString('ru')}</p>
                    )}
                </div>

                {/* Description */}
                <div className="bg-white rounded-2xl border border-primary-100 p-6 space-y-5">
                    <Section title="О вакансии" content={listing.description} />
                    <Section title="Требования" content={listing.requirements} />
                    <Section title="Мы предлагаем" content={listing.whatWeOffer} />
                </div>
            </div>

            {/* Sidebar */}
            <div className="space-y-4">
                {/* Apply card */}
                {user?.role === 'student' && (
                    <div className="bg-white rounded-2xl border border-primary-100 p-6">
                        <h3 className="font-semibold text-gray-900 mb-4">Откликнуться</h3>
                        {applied || applyMutation.isSuccess ? (
                            <div className="bg-green-50 text-green-700 rounded-xl px-4 py-3 text-sm font-medium">
                                ✓ Отклик отправлен
                            </div>
                        ) : (
                            <form onSubmit={handleSubmit((v) => applyMutation.mutate(v))} className="space-y-3">
                                <div>
                                    <label className="block text-xs font-medium text-gray-500 mb-1.5">
                                        Сопроводительное письмо
                                    </label>
                                    <textarea
                                        {...register('coverLetter', { required: 'Напишите пару слов о себе' })}
                                        rows={5}
                                        placeholder="Расскажите почему вы подходите на эту позицию..."
                                        className="w-full px-4 py-3 rounded-xl border border-gray-200 text-sm focus:outline-none focus:border-primary-400 focus:ring-2 focus:ring-primary-100 transition resize-none"
                                    />
                                    {errors.coverLetter && (
                                        <p className="text-xs text-red-500 mt-1">{errors.coverLetter.message}</p>
                                    )}
                                </div>
                                <button
                                    type="submit"
                                    disabled={applyMutation.isPending}
                                    className="w-full bg-primary-500 hover:bg-primary-600 text-white font-semibold py-3 rounded-xl transition-colors disabled:opacity-60"
                                >
                                    {applyMutation.isPending ? 'Отправляем...' : 'Откликнуться'}
                                </button>
                            </form>
                        )}
                    </div>
                )}

                {/* Company card */}
                {listing.company && (
                    <div className="bg-white rounded-2xl border border-primary-100 p-6">
                        <h3 className="font-semibold text-gray-900 mb-4">О компании</h3>
                        <div className="flex items-center gap-3 mb-3">
                            <div className="w-12 h-12 rounded-xl bg-primary-100 flex items-center justify-center text-primary-600 font-bold">
                                {listing.company.name[0]}
                            </div>
                            <div>
                                <p className="font-semibold text-gray-900">{listing.company.name}</p>
                                <p className="text-xs text-gray-500">{listing.company.industry}</p>
                            </div>
                        </div>
                        {listing.company.tagline && (
                            <p className="text-sm text-gray-600 mb-4">{listing.company.tagline}</p>
                        )}
                        <Link
                            to={`/companies/${listing.companyId}`}
                            className="block text-center border border-primary-200 text-primary-600 font-medium text-sm py-2.5 rounded-xl hover:bg-primary-50 transition-colors"
                        >
                            Профиль компании
                        </Link>
                    </div>
                )}
            </div>
        </div>
    )
}

function Section({ title, content }) {
    return (
        <div>
            <h2 className="font-semibold text-gray-900 mb-2">{title}</h2>
            <p className="text-sm text-gray-600 whitespace-pre-line leading-relaxed">{content}</p>
        </div>
    )
}