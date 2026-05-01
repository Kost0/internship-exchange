import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { profileApi } from '@/api/profile'
import SkillTag from '@/components/SkillTag'

const degreeLabels = { bachelor: 'Бакалавр', master: 'Магистр', phd: 'Аспирант' }
const levelLabels = { beginner: 'Начинающий', intermediate: 'Средний', advanced: 'Продвинутый' }

export default function StudentProfile() {
    const qc = useQueryClient()
    const [editingMain, setEditingMain] = useState(false)
    const [addingEdu, setAddingEdu] = useState(false)
    const [addingExp, setAddingExp] = useState(false)
    const [addingProject, setAddingProject] = useState(false)

    const { data: profile, isLoading } = useQuery({
        queryKey: ['student-profile'],
        queryFn: profileApi.getMyStudentProfile,
    })

    const updateMutation = useMutation({
        mutationFn: profileApi.updateStudentProfile,
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['student-profile'] })
            setEditingMain(false)
        },
    })

    const uploadAvatar = useMutation({
        mutationFn: profileApi.uploadAvatar,
        onSuccess: () => qc.invalidateQueries({ queryKey: ['student-profile'] }),
    })

    const uploadResume = useMutation({
        mutationFn: profileApi.uploadResume,
        onSuccess: () => qc.invalidateQueries({ queryKey: ['student-profile'] }),
    })

    const deleteEdu = useMutation({
        mutationFn: profileApi.deleteEducation,
        onSuccess: () => qc.invalidateQueries({ queryKey: ['student-profile'] }),
    })

    const deleteExp = useMutation({
        mutationFn: profileApi.deleteExperience,
        onSuccess: () => qc.invalidateQueries({ queryKey: ['student-profile'] }),
    })

    const deleteProject = useMutation({
        mutationFn: profileApi.deleteProject,
        onSuccess: () => qc.invalidateQueries({ queryKey: ['student-profile'] }),
    })

    const { register, handleSubmit, reset } = useForm({
        values: profile,
    })

    if (isLoading) return <div className="animate-pulse space-y-4"><div className="h-40 bg-white rounded-2xl" /></div>
    if (!profile) return null

    return (
        <div className="max-w-4xl mx-auto space-y-6">
            {/* Header card */}
            <div className="bg-white rounded-2xl border border-primary-100 p-6">
                <div className="flex items-start gap-5">
                    <div className="relative flex-shrink-0">
                        <div className="w-20 h-20 rounded-2xl bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-2xl overflow-hidden">
                            {profile.avatarUrl
                                ? <img src={profile.avatarUrl} className="w-full h-full object-cover" alt="avatar" />
                                : `${profile.firstName?.[0] ?? ''}${profile.lastName?.[0] ?? ''}`
                            }
                        </div>
                        <label className="absolute -bottom-1 -right-1 w-6 h-6 bg-primary-500 rounded-full flex items-center justify-center cursor-pointer hover:bg-primary-600 transition">
                            <span className="text-white text-xs">+</span>
                            <input type="file" accept="image/*" className="hidden"
                                   onChange={(e) => e.target.files?.[0] && uploadAvatar.mutate(e.target.files[0])}
                            />
                        </label>
                    </div>

                    {editingMain ? (
                        <form onSubmit={handleSubmit((v) => updateMutation.mutate(v))} className="flex-1 grid grid-cols-2 gap-3">
                            <input {...register('firstName')} placeholder="Имя" className="input-base" />
                            <input {...register('lastName')} placeholder="Фамилия" className="input-base" />
                            <input {...register('city')} placeholder="Город" className="input-base" />
                            <input {...register('phone')} placeholder="Телефон" className="input-base" />
                            <input {...register('githubUrl')} placeholder="GitHub URL" className="input-base" />
                            <input {...register('linkedinUrl')} placeholder="LinkedIn URL" className="input-base" />
                            <textarea {...register('bio')} placeholder="О себе" rows={2} className="input-base col-span-2 resize-none" />
                            <div className="col-span-2 flex gap-2">
                                <button type="submit" className="btn-primary text-sm px-4 py-2">Сохранить</button>
                                <button type="button" onClick={() => { setEditingMain(false); reset() }} className="btn-secondary text-sm px-4 py-2">Отмена</button>
                            </div>
                        </form>
                    ) : (
                        <div className="flex-1">
                            <div className="flex items-start justify-between">
                                <div>
                                    <h1 className="text-xl font-bold text-gray-900">{profile.firstName} {profile.lastName}</h1>
                                    <p className="text-sm text-gray-500 mt-0.5">{profile.city}</p>
                                    {profile.bio && <p className="text-sm text-gray-600 mt-2 max-w-lg">{profile.bio}</p>}
                                    <div className="flex gap-3 mt-3">
                                        {profile.githubUrl && <a href={profile.githubUrl} target="_blank" rel="noopener noreferrer" className="text-xs text-primary-600 hover:underline">GitHub</a>}
                                        {profile.linkedinUrl && <a href={profile.linkedinUrl} target="_blank" rel="noopener noreferrer" className="text-xs text-primary-600 hover:underline">LinkedIn</a>}
                                        {profile.portfolioUrl && <a href={profile.portfolioUrl} target="_blank" rel="noopener noreferrer" className="text-xs text-primary-600 hover:underline">Портфолио</a>}
                                    </div>
                                </div>
                                <div className="flex flex-col gap-2">
                                    <button onClick={() => setEditingMain(true)} className="btn-secondary text-sm px-4 py-2">Редактировать</button>
                                    <label className="btn-secondary text-sm px-4 py-2 cursor-pointer text-center">
                                        {profile.resumeUrl ? 'Заменить резюме' : 'Загрузить резюме'}
                                        <input type="file" accept=".pdf" className="hidden"
                                               onChange={(e) => e.target.files?.[0] && uploadResume.mutate(e.target.files[0])}
                                        />
                                    </label>
                                    {profile.resumeUrl && (
                                        <a href={profile.resumeUrl} target="_blank" rel="noopener noreferrer" className="text-xs text-primary-600 text-center hover:underline">
                                            Скачать PDF
                                        </a>
                                    )}
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Skills */}
            <div className="bg-white rounded-2xl border border-primary-100 p-6">
                <h2 className="font-semibold text-gray-900 mb-4">Навыки</h2>
                <div className="flex flex-wrap gap-2">
                    {profile.skills?.map((s) => (
                        <SkillTag key={s.id} label={`${s.skill} · ${levelLabels[s.level]}`} />
                    ))}
                </div>
            </div>

            {/* Languages */}
            <div className="bg-white rounded-2xl border border-primary-100 p-6">
                <h2 className="font-semibold text-gray-900 mb-4">Языки</h2>
                <div className="flex flex-wrap gap-2">
                    {profile.languages?.map((l) => (
                        <SkillTag key={l.id} label={`${l.language} · ${l.level}`} variant="gray" />
                    ))}
                </div>
            </div>

            {/* Education */}
            <div className="bg-white rounded-2xl border border-primary-100 p-6">
                <div className="flex items-center justify-between mb-4">
                    <h2 className="font-semibold text-gray-900">Образование</h2>
                    <button onClick={() => setAddingEdu(true)} className="text-sm text-primary-600 font-medium hover:underline">+ Добавить</button>
                </div>
                {addingEdu && (
                    <AddEducationForm
                        onSave={() => { qc.invalidateQueries({ queryKey: ['student-profile'] }); setAddingEdu(false) }}
                        onCancel={() => setAddingEdu(false)}
                    />
                )}
                <div className="space-y-4">
                    {profile.educations?.map((edu) => (
                        <EducationItem key={edu.id} edu={edu} onDelete={() => deleteEdu.mutate(edu.id)} />
                    ))}
                </div>
            </div>

            {/* Experience */}
            <div className="bg-white rounded-2xl border border-primary-100 p-6">
                <div className="flex items-center justify-between mb-4">
                    <h2 className="font-semibold text-gray-900">Опыт работы</h2>
                    <button onClick={() => setAddingExp(true)} className="text-sm text-primary-600 font-medium hover:underline">+ Добавить</button>
                </div>
                {addingExp && (
                    <AddExperienceForm
                        onSave={() => { qc.invalidateQueries({ queryKey: ['student-profile'] }); setAddingExp(false) }}
                        onCancel={() => setAddingExp(false)}
                    />
                )}
                <div className="space-y-4">
                    {profile.experiences?.map((exp) => (
                        <ExperienceItem key={exp.id} exp={exp} onDelete={() => deleteExp.mutate(exp.id)} />
                    ))}
                </div>
            </div>

            {/* Projects */}
            <div className="bg-white rounded-2xl border border-primary-100 p-6">
                <div className="flex items-center justify-between mb-4">
                    <h2 className="font-semibold text-gray-900">Проекты</h2>
                    <button onClick={() => setAddingProject(true)} className="text-sm text-primary-600 font-medium hover:underline">+ Добавить</button>
                </div>
                {addingProject && (
                    <AddProjectForm
                        onSave={() => { qc.invalidateQueries({ queryKey: ['student-profile'] }); setAddingProject(false) }}
                        onCancel={() => setAddingProject(false)}
                    />
                )}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    {profile.projects?.map((proj) => (
                        <ProjectItem key={proj.id} project={proj} onDelete={() => deleteProject.mutate(proj.id)} />
                    ))}
                </div>
            </div>
        </div>
    )
}

function EducationItem({ edu, onDelete }) {
    return (
        <div className="border-l-2 border-primary-200 pl-4">
            <div className="flex items-start justify-between">
                <div>
                    <p className="font-medium text-gray-900">{edu.university}</p>
                    <p className="text-sm text-gray-500">{edu.faculty} · {edu.specialization}</p>
                    <p className="text-sm text-gray-400">{degreeLabels[edu.degree]} · {edu.startYear} — {edu.isCurrent ? 'по н.в.' : edu.endYear}</p>
                    {edu.gpa && <p className="text-xs text-gray-400">GPA: {edu.gpa}</p>}
                </div>
                <button onClick={onDelete} className="text-xs text-red-400 hover:text-red-600 transition">Удалить</button>
            </div>
        </div>
    )
}

function ExperienceItem({ exp, onDelete }) {
    const formatLabels = { office: 'Офис', remote: 'Удалённо', hybrid: 'Гибрид' }
    return (
        <div className="border-l-2 border-primary-200 pl-4">
            <div className="flex items-start justify-between">
                <div>
                    <p className="font-medium text-gray-900">{exp.position}</p>
                    <p className="text-sm text-gray-500">{exp.companyName} · {formatLabels[exp.format]}</p>
                    <p className="text-sm text-gray-400">{exp.startDate} — {exp.isCurrent ? 'по н.в.' : exp.endDate}</p>
                    {exp.description && <p className="text-sm text-gray-600 mt-1">{exp.description}</p>}
                </div>
                <button onClick={onDelete} className="text-xs text-red-400 hover:text-red-600 transition">Удалить</button>
            </div>
        </div>
    )
}

function ProjectItem({ project, onDelete }) {
    return (
        <div className="border border-primary-100 rounded-xl p-4">
            <div className="flex items-start justify-between mb-1">
                <p className="font-medium text-gray-900">{project.title}</p>
                <button onClick={onDelete} className="text-xs text-red-400 hover:text-red-600 transition">Удалить</button>
            </div>
            {project.description && <p className="text-sm text-gray-500 mb-2">{project.description}</p>}
            <div className="flex flex-wrap gap-1.5 mb-2">
                {project.techs?.map((t) => <SkillTag key={t} label={t} />)}
            </div>
            {project.url && <a href={project.url} target="_blank" rel="noopener noreferrer" className="text-xs text-primary-600 hover:underline">Открыть проект →</a>}
        </div>
    )
}

function AddEducationForm({ onSave, onCancel }) {
    const { register, handleSubmit } = useForm()
    const mutation = useMutation({ mutationFn: profileApi.addEducation, onSuccess: onSave })
    return (
        <form onSubmit={handleSubmit((v) => mutation.mutate(v))} className="bg-primary-50 rounded-xl p-4 mb-4 space-y-3">
            <div className="grid grid-cols-2 gap-3">
                <input {...register('university', { required: true })} placeholder="Университет *" className="input-base" />
                <input {...register('faculty')} placeholder="Факультет" className="input-base" />
                <input {...register('specialization')} placeholder="Специальность" className="input-base" />
                <select {...register('degree')} className="input-base">
                    <option value="bachelor">Бакалавр</option>
                    <option value="master">Магистр</option>
                    <option value="phd">Аспирант</option>
                </select>
                <input {...register('startYear', { valueAsNumber: true })} type="number" placeholder="Год начала" className="input-base" />
                <input {...register('endYear', { valueAsNumber: true })} type="number" placeholder="Год окончания" className="input-base" />
            </div>
            <div className="flex gap-2">
                <button type="submit" className="btn-primary text-sm px-4 py-2">Сохранить</button>
                <button type="button" onClick={onCancel} className="btn-secondary text-sm px-4 py-2">Отмена</button>
            </div>
        </form>
    )
}

function AddExperienceForm({ onSave, onCancel }) {
    const { register, handleSubmit } = useForm()
    const mutation = useMutation({ mutationFn: profileApi.addExperience, onSuccess: onSave })
    return (
        <form onSubmit={handleSubmit((v) => mutation.mutate(v))} className="bg-primary-50 rounded-xl p-4 mb-4 space-y-3">
            <div className="grid grid-cols-2 gap-3">
                <input {...register('companyName', { required: true })} placeholder="Компания *" className="input-base" />
                <input {...register('position', { required: true })} placeholder="Должность *" className="input-base" />
                <input {...register('startDate')} type="date" className="input-base" />
                <input {...register('endDate')} type="date" className="input-base" />
                <select {...register('format')} className="input-base">
                    <option value="office">Офис</option>
                    <option value="remote">Удалённо</option>
                    <option value="hybrid">Гибрид</option>
                </select>
            </div>
            <textarea {...register('description')} placeholder="Описание" rows={2} className="input-base w-full resize-none" />
            <div className="flex gap-2">
                <button type="submit" className="btn-primary text-sm px-4 py-2">Сохранить</button>
                <button type="button" onClick={onCancel} className="btn-secondary text-sm px-4 py-2">Отмена</button>
            </div>
        </form>
    )
}

function AddProjectForm({ onSave, onCancel }) {
    const { register, handleSubmit } = useForm()
    const mutation = useMutation({
        mutationFn: (v) =>
            profileApi.addProject({ ...v, techs: v.techs.split(',').map((t) => t.trim()).filter(Boolean) }),
        onSuccess: onSave,
    })
    return (
        <form onSubmit={handleSubmit((v) => mutation.mutate(v))} className="bg-primary-50 rounded-xl p-4 mb-4 space-y-3">
            <input {...register('title', { required: true })} placeholder="Название проекта *" className="input-base w-full" />
            <input {...register('url')} placeholder="Ссылка (GitHub, деплой)" className="input-base w-full" />
            <input {...register('techs')} placeholder="Технологии через запятую: Go, React, Docker" className="input-base w-full" />
            <textarea {...register('description')} placeholder="Описание" rows={2} className="input-base w-full resize-none" />
            <div className="flex gap-2">
                <button type="submit" className="btn-primary text-sm px-4 py-2">Сохранить</button>
                <button type="button" onClick={onCancel} className="btn-secondary text-sm px-4 py-2">Отмена</button>
            </div>
        </form>
    )
}