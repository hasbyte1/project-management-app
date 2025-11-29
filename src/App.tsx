import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { RouterProvider } from '@tanstack/react-router';
import { useAuthStore } from './stores/auth-store';
import { Toaster } from 'sonner';
import { queryClient } from './lib/query-client';
import { router } from './router';

function App() {
  const auth = useAuthStore();

  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} context={{ auth, queryClient }} />
      <Toaster position="top-right" richColors />
      {import.meta.env.DEV && <ReactQueryDevtools />}
    </QueryClientProvider>
  );
}

export default App;
