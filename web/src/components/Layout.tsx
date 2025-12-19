import { Outlet, NavLink } from 'react-router-dom'

function Layout() {
  return (
    <div className="min-h-screen flex">
      <aside className="w-64 bg-dna-surface border-r border-dna-border p-4 hidden md:block">
        <h1 className="text-xl font-bold text-dna-primary mb-8">DNA Relations</h1>
        <nav className="space-y-2">
          <NavItem to="/">Home</NavItem>
          <NavItem to="/topics">Topics</NavItem>
          <NavItem to="/calendar">Calendar</NavItem>
          <NavItem to="/notes">Notes</NavItem>
          <NavItem to="/favorites">Favorites</NavItem>
          <NavItem to="/media">Media Sync</NavItem>
          <hr className="border-dna-border my-4" />
          <NavItem to="/settings">Settings</NavItem>
        </nav>
      </aside>
      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  )
}

function NavItem({ to, children }: { to: string; children: React.ReactNode }) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        `block px-4 py-2 rounded-lg transition-colors ${
          isActive
            ? 'bg-dna-primary text-white'
            : 'text-gray-400 hover:bg-dna-border hover:text-white'
        }`
      }
    >
      {children}
    </NavLink>
  )
}

export default Layout
