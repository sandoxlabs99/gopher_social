import { Outlet, createRootRoute } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'

// import { Header } from '@/components/Header'

export const API_URL =
  import.meta.env.VITE_API_URL || 'http://localhost:8080/v1'

export const Route = createRootRoute({
  component: () => (
    <>
      {/* <Header /> */}
      <Outlet />
      <TanStackDevtools
        config={{
          position: 'bottom-right',
        }}
        plugins={[
          {
            name: 'Tanstack Router',
            render: <TanStackRouterDevtoolsPanel />,
          },
        ]}
      />
    </>
  ),
})
