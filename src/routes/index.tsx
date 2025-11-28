import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/')({
  beforeLoad: ({ context }) => {
    if (!context.auth.isAuthenticated) {
      throw redirect({
        to: '/login',
      });
    }
    // Redirect authenticated users to their organizations
    throw redirect({
      to: '/organizations',
    });
  },
});
