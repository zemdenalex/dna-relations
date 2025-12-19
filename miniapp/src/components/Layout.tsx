import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { clearToken } from '../lib/api'

function Layout() {
  const navigate = useNavigate()

  const handleLogout = () => {
    clearToken()
    navigate('/login')
  }

  return (
    <div className="min-h-screen bg-dna-bg text-white flex flex-col">
      <main className="flex-1 p-4 pb-20 overflow-y-auto">
        <Outlet />
      </main>

      <nav className="fixed bottom-0 left-0 right-0 bg-dna-surface border-t border-dna-border">
        <div className="flex justify-around py-2">
          <NavItem to="/" icon="🏠" label="Home" />
          <NavItem to="/topics" icon="💬" label="Topics" />
          <NavItem to="/calendar" icon="📅" label="Calendar" />
          <NavItem to="/notes" icon="📝" label="Notes" />
          <button
            onClick={handleLogout}
            className="flex flex-col items-center py-1 px-3 text-gray-400 hover:text-red-400"
          >
            <span className="text-xl">🚪</span>
            <span className="text-xs mt-1">Logout</span>
          </button>
        </div>
      </nav>
    </div>
  )
}

function NavItem({
  to,
  icon,
  label,
}: {
  to: string
  icon: string
  label: string
}) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        `flex flex-col items-center py-1 px-3 ${
          isActive ? 'text-dna-primary' : 'text-gray-400'
        }`
      }
    >
      <span className="text-xl">{icon}</span>
      <span className="text-xs mt-1">{label}</span>
    </NavLink>
  )
}

export default Layout
