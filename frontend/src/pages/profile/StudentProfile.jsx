import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { profileApi } from '../../api/profile'

const degreeLabels = { bachelor: 'Бакалавр', master: 'Магистр', phd: 'Аспирант' }
const levelLabels = { beginner: 'Начинающий', intermediate: 'Средний', advanced: 'Продвинутый' }
const formatLabels = { office: 'Офис', remote: 'Удалённо', hybrid: 'Гибрид' }

export default function StudentProfile() {
    const qc = useQueryClient()
    const [editingMain, setEditingMain] = useState(false)
    const [addingEdu, setAddingEdu] = useState(false)
    const [addingExp, setAddingExp] = useState(false)
    const [addingProject, setAddingProject] = useState(false)
    const [addingSkill, setAddingSkill] = useState(false)
    const [addingLang, setAddingLang] = useState(false)

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

    const { register, handleSubmit, reset } = useForm({ values: profile })

    if (isLoading) return <div>Загрузка...</div>
    if (!profile) return null

    return (
        <div style={{ maxWidth: 800, margin: '0 auto' }}>

            {/* Header */}
            <div style={{ border: '1px solid #ccc', padding: 20, marginBottom: 16, borderRadius: 8 }}>
                <div style={{ display: 'flex', gap: 16, alignItems: 'flex-start' }}>
                    <div style={{ position: 'relative', flexShrink: 0 }}>
                        <div style={{ width: 72, height: 72, border: '1px solid #ccc', borderRadius: 8, overflow: 'hidden', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 22, fontWeight: 'bold', color: '#3e85dc', background: '#f5f5f5' }}>
                            {profile.avatar_url
                                ? <img src={profile.avatar_url} style={{ width: '100%', height: '100%', objectFit: 'cover' }} alt="avatar" />
                                : `${profile.first_name?.[0] ?? ''}${profile.last_name?.[0] ?? ''}`
                            }
                        </div>
                        <label style={{ position: 'absolute', bottom: -4, right: -4, width: 20, height: 20, background: '#3e85dc', color: 'white', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', fontSize: 14, borderRadius: 4 }}>
                            +
                            <input type="file" accept="image/*" style={{ display: 'none' }}
                                   onChange={(e) => e.target.files?.[0] && uploadAvatar.mutate(e.target.files[0])}
                            />
                        </label>
                    </div>

                    {editingMain ? (
                        <form onSubmit={handleSubmit((v) => updateMutation.mutate(v))} style={{ flex: 1 }}>
                            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8, marginBottom: 8 }}>
                                <input {...register('first_name')} placeholder="Имя" className="input-base" />
                                <input {...register('last_name')} placeholder="Фамилия" className="input-base" />
                                <input {...register('city')} placeholder="Город" className="input-base" />
                                <input {...register('phone')} placeholder="Телефон" className="input-base" />
                                <input {...register('github_url')} placeholder="GitHub URL" className="input-base" />
                                <input {...register('linkedin_url')} placeholder="LinkedIn URL" className="input-base" />
                            </div>
                            <textarea {...register('bio')} placeholder="О себе" rows={2} className="input-base" style={{ width: '100%', marginBottom: 8, resize: 'vertical' }} />
                            <div style={{ display: 'flex', gap: 8 }}>
                                <button type="submit" className="btn-primary">Сохранить</button>
                                <button type="button" className="btn-secondary" onClick={() => { setEditingMain(false); reset() }}>Отмена</button>
                            </div>
                        </form>
                    ) : (
                        <div style={{ flex: 1 }}>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                                <div>
                                    <div style={{ fontWeight: 'bold', fontSize: 18 }}>{profile.first_name} {profile.last_name}</div>
                                    {profile.city && <div style={{ fontSize: 13, color: '#555', marginTop: 2 }}>{profile.city}</div>}
                                    {profile.bio && <div style={{ fontSize: 13, marginTop: 6, maxWidth: 500 }}>{profile.bio}</div>}
                                    <div style={{ display: 'flex', gap: 12, marginTop: 8 }}>
                                        {profile.github_url && <a href={profile.github_url} target="_blank" rel="noopener noreferrer" style={{ fontSize: 13, color: '#3e85dc' }}>GitHub</a>}
                                        {profile.linkedin_url && <a href={profile.linkedin_url} target="_blank" rel="noopener noreferrer" style={{ fontSize: 13, color: '#3e85dc' }}>LinkedIn</a>}
                                        {profile.portfolio_url && <a href={profile.portfolio_url} target="_blank" rel="noopener noreferrer" style={{ fontSize: 13, color: '#3e85dc' }}>Портфолио</a>}
                                    </div>
                                </div>
                                <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                                    <button onClick={() => setEditingMain(true)} className="btn-secondary">Редактировать</button>
                                    <label className="btn-secondary" style={{ textAlign: 'center', cursor: 'pointer' }}>
                                        {profile.resume_url ? 'Заменить резюме' : 'Загрузить резюме'}
                                        <input type="file" accept=".pdf" style={{ display: 'none' }}
                                               onChange={(e) => e.target.files?.[0] && uploadResume.mutate(e.target.files[0])}
                                        />
                                    </label>
                                    {profile.resume_url && <ResumeLink studentId={profile.id} />}
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Skills */}
            <div style={{ border: '1px solid #ccc', padding: 20, marginBottom: 16, borderRadius: 8 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                    <div style={{ fontWeight: 'bold' }}>Навыки</div>
                    <button onClick={() => setAddingSkill(true)} style={{ fontSize: 13, color: '#3e85dc', background: 'none', border: 'none', cursor: 'pointer' }}>+ Добавить</button>
                </div>
                {addingSkill && (
                    <AddSkillForm
                        onSave={() => { qc.invalidateQueries({ queryKey: ['student-profile'] }); setAddingSkill(false) }}
                        onCancel={() => setAddingSkill(false)}
                    />
                )}
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
                    {profile.skills?.map(s => (
                        <span key={s.id} style={{ border: '1px solid #ccc', fontSize: 12, padding: '2px 8px', borderRadius: 4 }}>
                            {s.skill} · {levelLabels[s.level]}
                        </span>
                    ))}
                </div>
            </div>

            {/* Languages */}
            <div style={{ border: '1px solid #ccc', padding: 20, marginBottom: 16, borderRadius: 8 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                    <div style={{ fontWeight: 'bold' }}>Языки</div>
                    <button onClick={() => setAddingLang(true)} style={{ fontSize: 13, color: '#3e85dc', background: 'none', border: 'none', cursor: 'pointer' }}>+ Добавить</button>
                </div>
                {addingLang && (
                    <AddLanguageForm
                        onSave={() => { qc.invalidateQueries({ queryKey: ['student-profile'] }); setAddingLang(false) }}
                        onCancel={() => setAddingLang(false)}
                    />
                )}
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
                    {profile.languages?.map(l => (
                        <span key={l.id} style={{ border: '1px solid #ccc', fontSize: 12, padding: '2px 8px', borderRadius: 4 }}>
                            {l.language} · {l.level}
                        </span>
                    ))}
                </div>
            </div>

            {/* Education */}
            <div style={{ border: '1px solid #ccc', padding: 20, marginBottom: 16, borderRadius: 8 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                    <div style={{ fontWeight: 'bold' }}>Образование</div>
                    <button onClick={() => setAddingEdu(true)} style={{ fontSize: 13, color: '#3e85dc', background: 'none', border: 'none', cursor: 'pointer' }}>+ Добавить</button>
                </div>
                {addingEdu && (
                    <AddEducationForm
                        onSave={() => { qc.invalidateQueries({ queryKey: ['student-profile'] }); setAddingEdu(false) }}
                        onCancel={() => setAddingEdu(false)}
                    />
                )}
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    {profile.educations?.map(edu => (
                        <EducationItem key={edu.id} edu={edu} onDelete={() => deleteEdu.mutate(edu.id)} />
                    ))}
                </div>
            </div>

            {/* Experience */}
            <div style={{ border: '1px solid #ccc', padding: 20, marginBottom: 16, borderRadius: 8 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                    <div style={{ fontWeight: 'bold' }}>Опыт работы</div>
                    <button onClick={() => setAddingExp(true)} style={{ fontSize: 13, color: '#3e85dc', background: 'none', border: 'none', cursor: 'pointer' }}>+ Добавить</button>
                </div>
                {addingExp && (
                    <AddExperienceForm
                        onSave={() => { qc.invalidateQueries({ queryKey: ['student-profile'] }); setAddingExp(false) }}
                        onCancel={() => setAddingExp(false)}
                    />
                )}
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    {profile.experiences?.map(exp => (
                        <ExperienceItem key={exp.id} exp={exp} onDelete={() => deleteExp.mutate(exp.id)} />
                    ))}
                </div>
            </div>

            {/* Projects */}
            <div style={{ border: '1px solid #ccc', padding: 20, marginBottom: 16, borderRadius: 8 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                    <div style={{ fontWeight: 'bold' }}>Проекты</div>
                    <button onClick={() => setAddingProject(true)} style={{ fontSize: 13, color: '#3e85dc', background: 'none', border: 'none', cursor: 'pointer' }}>+ Добавить</button>
                </div>
                {addingProject && (
                    <AddProjectForm
                        onSave={() => { qc.invalidateQueries({ queryKey: ['student-profile'] }); setAddingProject(false) }}
                        onCancel={() => setAddingProject(false)}
                    />
                )}
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
                    {profile.projects?.map(proj => (
                        <ProjectItem key={proj.id} project={proj} onDelete={() => deleteProject.mutate(proj.id)} />
                    ))}
                </div>
            </div>
        </div>
    )
}

function ResumeLink({ studentId }) {
    const { data } = useQuery({
        queryKey: ['resume-link', studentId],
        queryFn: () => profileApi.getResumeUrl(studentId),
        enabled: !!studentId,
    })
    if (!data?.url) return null
    return (
        <a href={data.url} target="_blank" rel="noopener noreferrer" style={{ fontSize: 13, color: '#3e85dc', textAlign: 'center' }}>
            Скачать PDF
        </a>
    )
}

function EducationItem({ edu, onDelete }) {
    return (
        <div style={{ border: '1px solid #ddd', padding: 12, borderRadius: 6, display: 'flex', justifyContent: 'space-between' }}>
            <div>
                <div style={{ fontWeight: 'bold', fontSize: 14 }}>{edu.university}</div>
                {edu.faculty && <div style={{ fontSize: 13, color: '#555' }}>{edu.faculty} · {edu.specialization}</div>}
                <div style={{ fontSize: 13, color: '#777' }}>
                    {degreeLabels[edu.degree]} · {edu.start_year} — {edu.is_current ? 'по н.в.' : edu.end_year}
                </div>
                {edu.gpa > 0 && <div style={{ fontSize: 12, color: '#777' }}>GPA: {edu.gpa}</div>}
            </div>
            <button onClick={onDelete} style={{ fontSize: 12, color: '#e53e3e', background: 'none', border: 'none', cursor: 'pointer', alignSelf: 'flex-start' }}>
                Удалить
            </button>
        </div>
    )
}

function ExperienceItem({ exp, onDelete }) {
    return (
        <div style={{ border: '1px solid #ddd', padding: 12, borderRadius: 6, display: 'flex', justifyContent: 'space-between' }}>
            <div>
                <div style={{ fontWeight: 'bold', fontSize: 14 }}>{exp.position}</div>
                <div style={{ fontSize: 13, color: '#555' }}>{exp.company_name} · {formatLabels[exp.format]}</div>
                <div style={{ fontSize: 13, color: '#777' }}>
                    {exp.start_date} — {exp.is_current ? 'по н.в.' : exp.end_date}
                </div>
                {exp.description && <div style={{ fontSize: 13, marginTop: 4 }}>{exp.description}</div>}
            </div>
            <button onClick={onDelete} style={{ fontSize: 12, color: '#e53e3e', background: 'none', border: 'none', cursor: 'pointer', alignSelf: 'flex-start' }}>
                Удалить
            </button>
        </div>
    )
}

function ProjectItem({ project, onDelete }) {
    return (
        <div style={{ border: '1px solid #ddd', padding: 12, borderRadius: 6 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                <div style={{ fontWeight: 'bold', fontSize: 14 }}>{project.title}</div>
                <button onClick={onDelete} style={{ fontSize: 12, color: '#e53e3e', background: 'none', border: 'none', cursor: 'pointer' }}>
                    Удалить
                </button>
            </div>
            {project.description && <div style={{ fontSize: 13, color: '#555', marginBottom: 6 }}>{project.description}</div>}
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 4, marginBottom: 6 }}>
                {project.techs?.map(t => (
                    <span key={t} style={{ border: '1px solid #ccc', fontSize: 11, padding: '1px 6px', borderRadius: 4 }}>{t}</span>
                ))}
            </div>
            {project.url && (
                <a href={project.url} target="_blank" rel="noopener noreferrer" style={{ fontSize: 12, color: '#3e85dc' }}>
                    Открыть →
                </a>
            )}
        </div>
    )
}

function AddEducationForm({ onSave, onCancel }) {
    const { register, handleSubmit } = useForm()
    const mutation = useMutation({ mutationFn: profileApi.addEducation, onSuccess: onSave })
    return (
        <div style={{ border: '1px solid #ddd', padding: 12, marginBottom: 12, background: '#fafafa', borderRadius: 6 }}>
            <form onSubmit={handleSubmit((v) => mutation.mutate(v))}>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8, marginBottom: 8 }}>
                    <input {...register('university', { required: true })} placeholder="Университет *" className="input-base" />
                    <input {...register('faculty')} placeholder="Факультет" className="input-base" />
                    <input {...register('specialization')} placeholder="Специальность" className="input-base" />
                    <select {...register('degree')} className="input-base">
                        <option value="bachelor">Бакалавр</option>
                        <option value="master">Магистр</option>
                        <option value="phd">Аспирант</option>
                    </select>
                    <input {...register('start_year', { valueAsNumber: true })} type="number" placeholder="Год начала" className="input-base" />
                    <input {...register('end_year', { valueAsNumber: true })} type="number" placeholder="Год окончания" className="input-base" />
                </div>
                <div style={{ display: 'flex', gap: 8 }}>
                    <button type="submit" className="btn-primary">Сохранить</button>
                    <button type="button" className="btn-secondary" onClick={onCancel}>Отмена</button>
                </div>
            </form>
        </div>
    )
}

function AddExperienceForm({ onSave, onCancel }) {
    const { register, handleSubmit } = useForm()
    const mutation = useMutation({ mutationFn: profileApi.addExperience, onSuccess: onSave })
    return (
        <div style={{ border: '1px solid #ddd', padding: 12, marginBottom: 12, background: '#fafafa', borderRadius: 6 }}>
            <form onSubmit={handleSubmit((v) => mutation.mutate(v))}>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8, marginBottom: 8 }}>
                    <input {...register('company_name', { required: true })} placeholder="Компания *" className="input-base" />
                    <input {...register('position', { required: true })} placeholder="Должность *" className="input-base" />
                    <input {...register('start_date')} type="date" className="input-base" />
                    <input {...register('end_date')} type="date" className="input-base" />
                    <select {...register('format')} className="input-base">
                        <option value="office">Офис</option>
                        <option value="remote">Удалённо</option>
                        <option value="hybrid">Гибрид</option>
                    </select>
                </div>
                <textarea {...register('description')} placeholder="Описание" rows={2} className="input-base" style={{ width: '100%', marginBottom: 8, resize: 'vertical' }} />
                <div style={{ display: 'flex', gap: 8 }}>
                    <button type="submit" className="btn-primary">Сохранить</button>
                    <button type="button" className="btn-secondary" onClick={onCancel}>Отмена</button>
                </div>
            </form>
        </div>
    )
}

function AddProjectForm({ onSave, onCancel }) {
    const { register, handleSubmit } = useForm()
    const mutation = useMutation({
        mutationFn: (v) => profileApi.addProject({ ...v, techs: v.techs.split(',').map(t => t.trim()).filter(Boolean) }),
        onSuccess: onSave,
    })
    return (
        <div style={{ border: '1px solid #ddd', padding: 12, marginBottom: 12, background: '#fafafa', borderRadius: 6 }}>
            <form onSubmit={handleSubmit((v) => mutation.mutate(v))}>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8, marginBottom: 8 }}>
                    <input {...register('title', { required: true })} placeholder="Название *" className="input-base" />
                    <input {...register('url')} placeholder="Ссылка" className="input-base" />
                    <input {...register('techs')} placeholder="Технологии: Go, React, Docker" className="input-base" />
                    <textarea {...register('description')} placeholder="Описание" rows={2} className="input-base" style={{ resize: 'vertical' }} />
                </div>
                <div style={{ display: 'flex', gap: 8 }}>
                    <button type="submit" className="btn-primary">Сохранить</button>
                    <button type="button" className="btn-secondary" onClick={onCancel}>Отмена</button>
                </div>
            </form>
        </div>
    )
}

function AddSkillForm({ onSave, onCancel }) {
    const { register, handleSubmit } = useForm()
    const mutation = useMutation({ mutationFn: profileApi.addSkill, onSuccess: onSave })
    return (
        <div style={{ border: '1px solid #ddd', padding: 12, marginBottom: 12, background: '#fafafa', borderRadius: 6 }}>
            <form onSubmit={handleSubmit((v) => mutation.mutate(v))}>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8, marginBottom: 8 }}>
                    <input {...register('skill', { required: true })} placeholder="Навык *" className="input-base" />
                    <select {...register('level')} className="input-base">
                        <option value="beginner">Начинающий</option>
                        <option value="intermediate">Средний</option>
                        <option value="advanced">Продвинутый</option>
                    </select>
                </div>
                <div style={{ display: 'flex', gap: 8 }}>
                    <button type="submit" className="btn-primary">Сохранить</button>
                    <button type="button" className="btn-secondary" onClick={onCancel}>Отмена</button>
                </div>
            </form>
        </div>
    )
}

function AddLanguageForm({ onSave, onCancel }) {
    const { register, handleSubmit } = useForm()
    const mutation = useMutation({ mutationFn: profileApi.addLanguage, onSuccess: onSave })
    return (
        <div style={{ border: '1px solid #ddd', padding: 12, marginBottom: 12, background: '#fafafa', borderRadius: 6 }}>
            <form onSubmit={handleSubmit((v) => mutation.mutate(v))}>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8, marginBottom: 8 }}>
                    <input {...register('language', { required: true })} placeholder="Язык *" className="input-base" />
                    <select {...register('level')} className="input-base">
                        <option value="A1">A1</option>
                        <option value="A2">A2</option>
                        <option value="B1">B1</option>
                        <option value="B2">B2</option>
                        <option value="C1">C1</option>
                        <option value="C2">C2</option>
                        <option value="native">Родной</option>
                    </select>
                </div>
                <div style={{ display: 'flex', gap: 8 }}>
                    <button type="submit" className="btn-primary">Сохранить</button>
                    <button type="button" className="btn-secondary" onClick={onCancel}>Отмена</button>
                </div>
            </form>
        </div>
    )
}