import { Outlet } from 'react-router-dom'
import Navbar from './Navbar'

export default function Layout() {
    return (
        <div style={{ minHeight: '100vh', background: 'white' }}>
            <Navbar />
            <div style={{ maxWidth: 1100, margin: '0 auto', padding: '24px 16px' }}>
                <Outlet />
            </div>
        </div>
    )
}