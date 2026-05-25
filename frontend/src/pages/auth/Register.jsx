import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { authApi } from '../../api/auth'
import { useAuthStore } from '../../store/auth'

export default function Register() {
    const navigate = useNavigate()
    const { setAuth } = useAuthStore()
    const [role, setRole] = useState('student')
    const { register, handleSubmit, watch, formState: { errors } } = useForm()

    const mutation = useMutation({
        mutationFn: (v) => authApi.register({ email: v.email, password: v.password, role }),
        onSuccess: (data) => {
            setAuth(data.user, data.tokens.accessToken, data.tokens.refreshToken)
            navigate('/dashboard')
        },
    })

    return (
        <div style={{ maxWidth: 400, margin: '60px auto', border: '1px solid #ccc', padding: 24, borderRadius: 10 }}>
            <h2 style={{ marginTop: 0, marginBottom: 20 }}>Регистрация</h2>

            <div style={{ marginBottom: 16, display: 'flex', gap: 0 }}>
                <button
                    type="button"
                    onClick={() => setRole('student')}
                    style={{ flex: 1, padding: '8px', border: '1px solid #ccc', background: role === 'student' ? '#3e85dc' : 'white', color: role === 'student' ? 'white' : '#333', cursor: 'pointer', fontSize: 13 }}
                >
                    Студент
                </button>
                <button
                    type="button"
                    onClick={() => setRole('company')}
                    style={{ flex: 1, padding: '8px', border: '1px solid #ccc', borderLeft: 'none', background: role === 'company' ? '#3e85dc' : 'white', color: role === 'company' ? 'white' : '#333', cursor: 'pointer', fontSize: 13 }}
                >
                    Компания
                </button>
            </div>

            <form onSubmit={handleSubmit((v) => mutation.mutate(v))}>
                <div style={{ marginBottom: 12 }}>
                    <div style={{ fontSize: 13, marginBottom: 4 }}>Email</div>
                    <input {...register('email', { required: true })} type="email" className="input-base" />
                </div>
                <div style={{ marginBottom: 12 }}>
                    <div style={{ fontSize: 13, marginBottom: 4 }}>Пароль</div>
                    <input {...register('password', { required: true, minLength: 8 })} type="password" className="input-base" />
                </div>
                <div style={{ marginBottom: 16 }}>
                    <div style={{ fontSize: 13, marginBottom: 4 }}>Повторите пароль</div>
                    <input {...register('confirmPassword', { validate: v => v === watch('password') || 'Пароли не совпадают' })} type="password" className="input-base" />
                    {errors.confirmPassword && <div style={{ color: 'red', fontSize: 12, marginTop: 4 }}>{errors.confirmPassword.message}</div>}
                </div>
                {mutation.isError && <div style={{ color: 'red', fontSize: 13, marginBottom: 12 }}>Ошибка регистрации</div>}
                <button type="submit" className="btn-primary" style={{ width: '100%' }}>
                    {mutation.isPending ? 'Создаём...' : 'Зарегистрироваться'}
                </button>
            </form>
            <div style={{ marginTop: 16, fontSize: 13 }}>
                Уже есть аккаунт? <Link to="/login">Войти</Link>
            </div>
        </div>
    )
}