import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/store/auth'

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
        <div className="min-h-screen bg-primary-50 flex items-center justify-center px-4">
            <div className="bg-white rounded-3xl border border-primary-100 p-10 w-full max-w-md shadow-sm">
                <h1 className="text-2xl font-bold text-gray-900 mb-1">Добро пожаловать</h1>
                <p className="text-sm text-gray-500 mb-8">
                    Нет аккаунта?{' '}
                    <Link to="/register" className="text-primary-600 font-medium hover:underline">
                        Зарегистрироваться
                    </Link>
                </p>

                <form onSubmit={handleSubmit((v) => mutation.mutate(v))} className="space-y-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1.5">Email</label>
                        <input
                            {...register('email', { required: 'Обязательное поле' })}
                            type="email"
                            placeholder="ivan@university.ru"
                            className="w-full px-4 py-3 rounded-xl border border-gray-200 text-sm focus:outline-none focus:border-primary-400 focus:ring-2 focus:ring-primary-100 transition"
                        />
                        {errors.email && <p className="text-xs text-red-500 mt-1">{errors.email.message}</p>}
                    </div>

                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1.5">Пароль</label>
                        <input
                            {...register('password', {
                                required: 'Обязательное поле',
                                minLength: { value: 8, message: 'Минимум 8 символов' }
                            })}
                            type="password"
                            placeholder="••••••••"
                            className="w-full px-4 py-3 rounded-xl border border-gray-200 text-sm focus:outline-none focus:border-primary-400 focus:ring-2 focus:ring-primary-100 transition"
                        />
                        {errors.password && <p className="text-xs text-red-500 mt-1">{errors.password.message}</p>}
                    </div>

                    {mutation.isError && (
                        <p className="text-sm text-red-500 bg-red-50 rounded-xl px-4 py-3">
                            Неверный email или пароль
                        </p>
                    )}

                    <button
                        type="submit"
                        disabled={mutation.isPending}
                        className="w-full bg-primary-500 hover:bg-primary-600 text-white font-semibold py-3 rounded-xl transition-colors disabled:opacity-60 mt-2"
                    >
                        {mutation.isPending ? 'Входим...' : 'Войти'}
                    </button>
                </form>
            </div>
        </div>
    )
}