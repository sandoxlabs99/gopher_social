import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  return (
    <main className="bg-background text-foreground flex flex-col">
      <div>This is the content of the root page</div>
    </main>
  )
}
