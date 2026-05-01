import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/store/auth'
import clsx from 'clsx'

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
        <div className="min-h-screen bg-primary-50 flex items-center justify-center px-4 py-12">
            <div className="bg-white rounded-3xl border border-primary-100 p-10 w-full max-w-md shadow-sm">
                <h1 className="text-2xl font-bold text-gray-900 mb-1">Создать аккаунт</h1>
                <p className="text-sm text-gray-500 mb-6">
                    Уже есть аккаунт?{' '}
                    <Link to="/login" className="text-primary-600 font-medium hover:underline">
                        Войти
                    </Link>
                </p>

                <div className="flex bg-gray-100 rounded-xl p-1 mb-6">
                    {['student', 'company'].map((r) => (
                        <button
                            key={r}
                            type="button"
                            onClick={() => setRole(r)}
                            className={clsx(
                                'flex-1 py-2 text-sm font-medium rounded-lg transition-all',
                                role === r ? 'bg-white text-primary-600 shadow-sm' : 'text-gray-500'
                            )}
                        >
                            {r === 'student' ? 'Я студент' : 'Я компания'}
                        </button>
                    ))}
                </div>

                <form onSubmit={handleSubmit((v) => mutation.mutate(v))} className="space-y-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1.5">Email</label>
                        <input
                            {...register('email', { required: 'Обязательное поле' })}
                            type="email"
                            placeholder={role === 'student' ? 'ivan@university.ru' : 'hr@company.ru'}
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

                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1.5">Повторите пароль</label>
                        <input
                            {...register('confirmPassword', {
                                required: 'Обязательное поле',
                                validate: (v) => v === watch('password') || 'Пароли не совпадают',
                            })}
                            type="password"
                            placeholder="••••••••"
                            className="w-full px-4 py-3 rounded-xl border border-gray-200 text-sm focus:outline-none focus:border-primary-400 focus:ring-2 focus:ring-primary-100 transition"
                        />
                        {errors.confirmPassword && <p className="text-xs text-red-500 mt-1">{errors.confirmPassword.message}</p>}
                    </div>

                    {mutation.isError && (
                        <p className="text-sm text-red-500 bg-red-50 rounded-xl px-4 py-3">
                            Ошибка регистрации. Попробуйте другой email.
                        </p>
                    )}

                    <button
                        type="submit"
                        disabled={mutation.isPending}
                        className="w-full bg-primary-500 hover:bg-primary-600 text-white font-semibold py-3 rounded-xl transition-colors disabled:opacity-60 mt-2"
                    >
                        {mutation.isPending ? 'Создаём аккаунт...' : 'Зарегистрироваться'}
                    </button>
                </form>
            </div>
        </div>
    )
}