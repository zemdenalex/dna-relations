import { Routes, Route, Navigate } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { isAuthenticated } from './lib/api'
import Layout from './components/Layout'
import HomePage from './pages/Home'
import TopicsPage from './pages/Topics'
import CalendarPage from './pages/Calendar'
import NotesPage from './pages/Notes'
import LoginPage from './pages/Login'

declare global {
  interface Window {
    Telegram: {
      WebApp: {
        ready: () => void
        expand: () => void
        MainButton: {
          show: () => void
          hide: () => void
          setText: (text: string) => void
          onClick: (fn: () => void) => void
        }
        BackButton: {
          show: () => void
          hide: () => void
          onClick: (fn: () => void) => void
        }
        themeParams: Record<string, string>
        initDataUnsafe: {
          user?: {
            id: number
            first_name: string
            last_name?: string
            username?: string
          }
        }
      }
    }
  }
}

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}

function App() {
  const [ready, setReady] = useState(false)

  useEffect(() => {
    if (window.Telegram?.WebApp) {
      window.Telegram.WebApp.ready()
      window.Telegram.WebApp.expand()
    }
    setReady(true)
  }, [])

  if (!ready) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-dna-primary">Loading...</div>
      </div>
    )
  }

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <Layout />
          </ProtectedRoute>
        }
      >
        <Route index element={<HomePage />} />
        <Route path="topics" element={<TopicsPage />} />
        <Route path="calendar" element={<CalendarPage />} />
        <Route path="notes" element={<NotesPage />} />
      </Route>
    </Routes>
  )
}

export default App
