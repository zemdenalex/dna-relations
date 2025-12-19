function HomePage() {
  return (
    <div className="space-y-6">
      <header>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-gray-400">Welcome to DNA Relations</p>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <StatCard title="Topics" value="--" subtitle="pending" />
        <StatCard title="Events" value="--" subtitle="this week" />
        <StatCard title="Notes" value="--" subtitle="total" />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <section className="bg-dna-surface rounded-lg p-6">
          <h2 className="font-semibold mb-4">Priority Topics</h2>
          <p className="text-gray-400">[stub] No topics</p>
        </section>

        <section className="bg-dna-surface rounded-lg p-6">
          <h2 className="font-semibold mb-4">Upcoming Events</h2>
          <p className="text-gray-400">[stub] No events</p>
        </section>
      </div>

      <section className="bg-dna-surface rounded-lg p-6">
        <h2 className="font-semibold mb-4">Pinned Notes</h2>
        <p className="text-gray-400">[stub] No pinned notes</p>
      </section>
    </div>
  )
}

function StatCard({
  title,
  value,
  subtitle,
}: {
  title: string
  value: string
  subtitle: string
}) {
  return (
    <div className="bg-dna-surface rounded-lg p-6">
      <h3 className="text-gray-400 text-sm">{title}</h3>
      <div className="text-3xl font-bold text-dna-primary mt-2">{value}</div>
      <p className="text-gray-500 text-sm mt-1">{subtitle}</p>
    </div>
  )
}

export default HomePage
