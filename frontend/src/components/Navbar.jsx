import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '../store/auth'

function isTokenExpired(token) {
    if (!token) return true
    try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        return payload.exp * 1000 < Date.now()
    } catch { return true }
}

export default function Navbar() {
    const { user, accessToken, clearAuth } = useAuthStore()
    const navigate = useNavigate()
    const location = useLocation()
    const isAuth = user && !isTokenExpired(accessToken)

    const handleLogout = () => { clearAuth(); navigate('/login') }

    return (
        <div style={{ borderBottom: '1px solid #ccc', padding: '10px 24px', display: 'flex', alignItems: 'center', justifyContent: 'space-between', borderRadius: 6 }}>
            <Link to="/" style={{ fontWeight: 'bold', fontSize: 18, color: '#3e85dc', textDecoration: 'none' }}>
                Биржа стажировок
            </Link>
            <div style={{ display: 'flex', gap: 20, alignItems: 'center' }}>
                <Link to="/" style={{ color: location.pathname === '/' ? '#3e85dc' : '#333', textDecoration: 'none', fontSize: 14 }}>
                    Вакансии
                </Link>
                {isAuth ? (
                    <>
                        <Link to="/dashboard" style={{ color: '#333', textDecoration: 'none', fontSize: 14 }}>Кабинет</Link>
                        <Link to={user.role === 'student' ? '/profile' : '/company/profile'} style={{ color: '#333', textDecoration: 'none', fontSize: 14 }}>Профиль</Link>
                        <button onClick={handleLogout} style={{ fontSize: 14, background: 'none', border: 'none', cursor: 'pointer', color: '#333' }}>Выйти</button>
                    </>
                ) : (
                    <>
                        <Link to="/login" style={{ color: '#333', textDecoration: 'none', fontSize: 14 }}>Войти</Link>
                        <Link to="/register" style={{ background: '#3e85dc', color: 'white', padding: '6px 14px', textDecoration: 'none', fontSize: 14 }}>Регистрация</Link>
                    </>
                )}
            </div>
        </div>
    )
}