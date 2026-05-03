import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/auth'
import Layout from './components/Layout'
import Login from './pages/auth/Login'
import Register from './pages/auth/Register'
import ListingsPage from './pages/listings/ListingPage'
import ListingDetail from './pages/listings/ListingDetail'
import StudentProfile from './pages/profile/StudentProfile'
import CompanyProfile from './pages/profile/CompanyProfile'
import StudentDashboard from './pages/dashboard/StudentDashboard'
import CompanyDashboard from './pages/dashboard/CompanyDashboard'

function RequireAuth({ children }) {
  const user = useAuthStore((s) => s.user)
  return user ? <>{children}</> : <Navigate to="/login" replace />
}

function RequireRole({ role, children }) {
  const user = useAuthStore((s) => s.user)
  if (!user) return <Navigate to="/login" replace />
  if (user.role !== role) return <Navigate to="/" replace />
  return <>{children}</>
}

export default function App() {
  return (
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route element={<Layout />}>
          <Route path="/" element={<ListingsPage />} />
          <Route path="/listings/:id" element={<ListingDetail />} />
          <Route path="/companies/:id" element={<CompanyProfile />} />

          <Route
              path="/profile"
              element={
                <RequireRole role="student">
                  <StudentProfile />
                </RequireRole>
              }
          />
          <Route
              path="/company/profile"
              element={
                <RequireRole role="company">
                  <CompanyProfile own />
                </RequireRole>
              }
          />
          <Route
              path="/dashboard"
              element={
                <RequireAuth>
                  <DashboardRedirect />
                </RequireAuth>
              }
          />
          <Route
              path="/dashboard/student"
              element={
                <RequireRole role="student">
                  <StudentDashboard />
                </RequireRole>
              }
          />
          <Route
              path="/dashboard/company"
              element={
                <RequireRole role="company">
                  <CompanyDashboard />
                </RequireRole>
              }
          />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
  )
}

function DashboardRedirect() {
  const user = useAuthStore((s) => s.user)
  if (user?.role === 'student') return <Navigate to="/dashboard/student" replace />
  if (user?.role === 'company') return <Navigate to="/dashboard/company" replace />
  return <Navigate to="/" replace />
}