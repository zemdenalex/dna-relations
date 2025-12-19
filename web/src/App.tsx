import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import HomePage from './pages/Home'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="topics" element={<Stub name="Topics" />} />
        <Route path="calendar" element={<Stub name="Calendar" />} />
        <Route path="notes" element={<Stub name="Notes" />} />
        <Route path="favorites" element={<Stub name="Favorites" />} />
        <Route path="media" element={<Stub name="Media Sync" />} />
        <Route path="settings" element={<Stub name="Settings" />} />
        <Route path="login" element={<Stub name="Login" />} />
      </Route>
    </Routes>
  )
}

function Stub({ name }: { name: string }) {
  return (
    <div className="flex items-center justify-center h-64">
      <div className="text-center">
        <h1 className="text-2xl font-bold mb-2">{name}</h1>
        <p className="text-gray-400">[stub] Not implemented</p>
      </div>
    </div>
  )
}

export default App
