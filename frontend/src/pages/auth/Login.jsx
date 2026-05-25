import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { authApi } from '../../api/auth'
import { useAuthStore } from '../../store/auth'

export default function Login() {
    const navigate = useNavigate()
    const { setAuth } = useAuthStore()
    const { register, handleSubmit, formState: { errors } } = useForm()

    const mutation = useMutation({
        mutationFn: authApi.login,
        onSuccess: (data) => {
            setAuth(data.user, data.tokens.accessToken, data.tokens.refreshToken)
            navigate('/dashboard')
        },
    })

    return (
        <div style={{ maxWidth: 400, margin: '60px auto', border: '1px solid #ccc', padding: 24, borderRadius: 10 }}>
            <h2 style={{ marginTop: 0, marginBottom: 20 }}>Войти</h2>

            <form onSubmit={handleSubmit((v) => mutation.mutate(v))}>
                <div style={{ marginBottom: 12 }}>
                    <div style={{ fontSize: 13, marginBottom: 4 }}>Email</div>
                    <input {...register('email', { required: true })} type="email" className="input-base" />
                </div>
                <div style={{ marginBottom: 16 }}>
                    <div style={{ fontSize: 13, marginBottom: 4 }}>Пароль</div>
                    <input {...register('password', { required: true })} type="password" className="input-base" />
                </div>
                {mutation.isError && <div style={{ color: 'red', fontSize: 13, marginBottom: 12 }}>Неверный email или пароль</div>}
                <button type="submit" className="btn-primary" style={{ width: '100%' }}>
                    {mutation.isPending ? 'Входим...' : 'Войти'}
                </button>
            </form>
            <div style={{ marginTop: 16, fontSize: 13 }}>
                Нет аккаунта? <Link to="/register">Зарегистрироваться</Link>
            </div>
        </div>
    )
}