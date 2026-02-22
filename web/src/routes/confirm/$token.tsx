import { createFileRoute, Link, notFound } from '@tanstack/react-router'
import { API_URL } from '@/routes/__root'

export const Route = createFileRoute('/confirm/$token')({
  component: RouteComponent,
  loader: async ({ params }) => {
    const response = await fetch(`${API_URL}/users/activate/${params.token}`, {
      method: 'PUT',
    })

    if (response.status === 404) {
      throw notFound()
    }

    if (!response.ok) {
      throw new Error('Failed to activate account. Please try again later.')
    }

    throw Route.redirect({
      to: '/',
    })
  },
  pendingComponent: () => <div className="px-4 py-6">Loading...</div>,
  errorComponent: () => (
    <div className="px-4 py-6">
      An error occured! Failed to activate account. Please try again later.
    </div>
  ),
  notFoundComponent: () => <div className="px-4 py-6">Invalid token</div>,
})

function RouteComponent() {
  // const { id: paramId } = Route.useParams() // alternative way to access params

  return (
    <div className="px-4 py-6">
      Your account has been activated! You'll be redirected shortly. If not, click <Link to="/">here</Link>.
    </div>
  )
}
