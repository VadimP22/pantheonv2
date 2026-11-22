import { QueryClient } from '@tanstack/react-query'

/** Единственный экземпляр QueryClient для всего приложения. */
export const queryClient = new QueryClient()