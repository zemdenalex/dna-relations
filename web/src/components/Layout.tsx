import Sidebar from './Sidebar'

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-gray-900">
      <Sidebar />
      <main className="ml-56 p-8">{children}</main>
    </div>
  )
}