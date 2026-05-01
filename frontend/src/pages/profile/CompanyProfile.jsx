import { useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { useState } from 'react'
import { profileApi } from '@/api/profile'
import { listingsApi } from '@/api/listings'
import JobCard from '@/components/JobCard'

export default function CompanyProfile({ own = false }) {
    const { id } = useParams()
    const qc = useQueryClient()
    const [editing, setEditing] = useState(false)

    const { data: profile, isLoading } = useQuery({
        queryKey: ['company-profile', own ? 'own' : id],
        queryFn: () => own ? profileApi.getMyCompanyProfile() : profileApi.getCompanyProfile(id),
        enabled: own || !!id,
    })

    const { data: listings } = useQuery({
        queryKey: ['company-listings', profile?.id],
        queryFn: () => listingsApi.getListings({ page: 1, limit: 10, companyId: profile?.id }),
        enabled: !!profile,
    })

    const updateMutation = useMutation({
        mutationFn: profileApi.updateCompanyProfile,
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['company-profile'] })
            setEditing(false)
        },
    })

    const uploadLogo = useMutation({
        mutationFn: profileApi.uploadLogo,
        onSuccess: () => qc.invalidateQueries({ queryKey: ['company-profile'] }),
    })

    const { register, handleSubmit } = useForm({ values: profile })

    if (isLoading) return <div className="animate-pulse h-64 bg-white rounded-2xl" />
    if (!profile) return <p className="text-center text-gray-500 py-20">Компания не найдена</p>

    return (
        <div className="max-w-5xl mx-auto space-y-6">
            <div className="bg-primary-100 rounded-3xl p-8">
                <div className="flex items-end gap-5">
                    <div className="relative">
                        <div className="w-24 h-24 rounded-2xl bg-primary-500 flex items-center justify-center text-white font-bold text-3xl overflow-hidden shadow-lg">
                            {profile.logoUrl
                                ? <img src={profile.logoUrl} className="w-full h-full object-cover" alt="logo" />
                                : profile.name?.[0]
                            }
                        </div>
                        {own && (
                            <label className="absolute -bottom-1 -right-1 w-7 h-7 bg-primary-500 rounded-full flex items-center justify-center cursor-pointer hover:bg-primary-600 transition shadow">
                                <span className="text-white text-sm">+</span>
                                <input type="file" accept="image/*" className="hidden"
                                       onChange={(e) => e.target.files?.[0] && uploadLogo.mutate(e.target.files[0])}
                                />
                            </label>
                        )}
                    </div>
                    <div className="flex-1">
                        <h1 className="text-2xl font-bold text-gray-900">{profile.name}</h1>
                        <p className="text-gray-500">
                            {profile.industry} {profile.size ? `· ${profile.size} сотрудников` : ''}
                        </p>
                        {profile.website && (
                            <a href={profile.website} target="_blank" rel="noopener noreferrer" className="text-sm text-primary-600 hover:underline">
                                {profile.website}
                            </a>
                        )}
                    </div>
                    {own && (
                        <button onClick={() => setEditing(!editing)} className="btn-secondary text-sm px-5 py-2.5">
                            {editing ? 'Отмена' : 'Редактировать'}
                        </button>
                    )}
                </div>
            </div>

            {editing && (
                <form onSubmit={handleSubmit((v) => updateMutation.mutate(v))} className="bg-white rounded-2xl border border-primary-100 p-6 space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <input {...register('name')} placeholder="Название компании" className="input-base" />
                        <input {...register('tagline')} placeholder="Короткое описание" className="input-base" />
                        <input {...register('industry')} placeholder="Сфера деятельности" className="input-base" />
                        <input {...register('size')} placeholder="Размер компании" className="input-base" />
                        <input {...register('city')} placeholder="Город" className="input-base" />
                        <input {...register('website')} placeholder="Сайт" className="input-base" />
                        <input {...register('contactEmail')} placeholder="Email для связи" className="input-base" />
                        <input {...register('foundedYear', { valueAsNumber: true })} type="number" placeholder="Год основания" className="input-base" />
                    </div>
                    <textarea {...register('description')} rows={4} placeholder="Полное описание компании" className="input-base w-full resize-none" />
                    <button type="submit" className="btn-primary px-6 py-2.5">Сохранить</button>
                </form>
            )}

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="bg-white rounded-2xl border border-primary-100 p-5 md:col-span-2">
                    <h2 className="font-semibold text-gray-900 mb-3">О компании</h2>
                    <p className="text-sm text-gray-600 leading-relaxed">{profile.description ?? 'Описание не указано'}</p>
                    {profile.tagline && <p className="text-sm font-medium text-primary-600 mt-3">{profile.tagline}</p>}
                </div>
                <div className="bg-white rounded-2xl border border-primary-100 p-5">
                    <h2 className="font-semibold text-gray-900 mb-3">Контакты</h2>
                    <div className="space-y-2 text-sm">
                        {profile.city && <Row label="Город" value={profile.city} />}
                        {profile.contactEmail && <Row label="Email" value={profile.contactEmail} />}
                        {profile.foundedYear && <Row label="Основана" value={String(profile.foundedYear)} />}
                        <Row label="Формат" value={profile.isRemoteFriendly ? 'Удалёнка ок' : 'Только офис'} />
                    </div>
                </div>
            </div>

            {listings?.items && listings.items.length > 0 && (
                <div>
                    <h2 className="font-semibold text-gray-900 mb-4">Открытые вакансии</h2>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                        {listings.items.map((l) => <JobCard key={l.id} listing={l} />)}
                    </div>
                </div>
            )}
        </div>
    )
}

function Row({ label, value }) {
    return (
        <div className="flex justify-between">
            <span className="text-gray-400">{label}</span>
            <span className="text-gray-700 font-medium">{value}</span>
        </div>
    )
}