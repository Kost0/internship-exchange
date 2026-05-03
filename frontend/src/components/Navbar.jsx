import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '../store/auth'
import clsx from 'clsx'

export default function Navbar() {
    const { user, clearAuth } = useAuthStore()
    const navigate = useNavigate()
    const location = useLocation()

    const handleLogout = () => {
        clearAuth()
        navigate('/login')
    }

    const isActive = (path) => location.pathname === path

    return (
        <nav className="bg-white border-b border-primary-100 sticky top-0 z-50">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <Link to="/" className="text-xl font-bold text-primary-600">
                        Стажировки
                    </Link>

                    <div className="flex items-center gap-6">
                        <Link
                            to="/"
                            className={clsx(
                                'text-sm font-medium transition-colors',
                                isActive('/') ? 'text-primary-600' : 'text-gray-500 hover:text-gray-900'
                            )}
                        >
                            Вакансии
                        </Link>

                        {user ? (
                            <>
                                <Link
                                    to="/dashboard"
                                    className={clsx(
                                        'text-sm font-medium transition-colors',
                                        location.pathname.startsWith('/dashboard')
                                            ? 'text-primary-600'
                                            : 'text-gray-500 hover:text-gray-900'
                                    )}
                                >
                                    Кабинет
                                </Link>
                                <Link
                                    to={user.role === 'student' ? '/profile' : '/company/profile'}
                                    className={clsx(
                                        'text-sm font-medium transition-colors',
                                        location.pathname.includes('profile')
                                            ? 'text-primary-600'
                                            : 'text-gray-500 hover:text-gray-900'
                                    )}
                                >
                                    Профиль
                                </Link>
                                <button
                                    onClick={handleLogout}
                                    className="text-sm font-medium text-gray-500 hover:text-gray-900 transition-colors"
                                >
                                    Выйти
                                </button>
                            </>
                        ) : (
                            <>
                                <Link
                                    to="/login"
                                    className="text-sm font-medium text-gray-500 hover:text-gray-900"
                                >
                                    Войти
                                </Link>
                                <Link
                                    to="/register"
                                    className="bg-primary-500 text-white text-sm font-medium px-4 py-2 rounded-xl hover:bg-primary-600 transition-colors"
                                >
                                    Регистрация
                                </Link>
                            </>
                        )}
                    </div>
                </div>
            </div>
        </nav>
    )
}