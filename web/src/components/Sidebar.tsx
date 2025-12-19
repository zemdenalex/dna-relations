import { NavLink } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

export default function Sidebar() {
  const { user, logout } = useAuth()

  const links = [
    { to: '/', label: 'Home' },
    { to: '/topics', label: 'Topics' },
    { to: '/calendar', label: 'Calendar' },
    { to: '/notes', label: 'Notes' },
  ]

  return (
    <div className="w-56 bg-gray-900 h-screen fixed left-0 top-0 flex flex-col">
      <div className="p-6">
        <h1 className="text-xl font-bold text-indigo-400">DNA Relations</h1>
      </div>

      <nav className="flex-1 px-4">
        {links.map(link => (
          <NavLink
            key={link.to}
            to={link.to}
            end={link.to === '/'}
            className={({ isActive }) =>
              `block px-4 py-3 rounded-lg mb-1 transition-colors ${
                isActive ? 'bg-indigo-600 text-white' : 'text-gray-300 hover:bg-gray-800'
              }`
            }
          >
            {link.label}
          </NavLink>
        ))}
      </nav>

      <div className="p-4 border-t border-gray-800">
        <div className="text-gray-400 text-sm mb-2">{user?.display_name || user?.username}</div>
        <button onClick={logout} className="text-gray-500 hover:text-gray-300 text-sm">
          Logout
        </button>
      </div>
    </div>
  )
}