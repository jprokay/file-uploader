'use client'

import { HydrationBoundary, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import * as React from 'react'

export function Providers(props: { children: React.ReactNode }) {
  const [queryClient] = React.useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 5 * 1000,
          },
        },
      }),
  )

  return (
    <QueryClientProvider client={queryClient}>
      <HydrationBoundary>
        {props.children}
      </HydrationBoundary>
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  )
}
