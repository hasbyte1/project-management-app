import { createFileRoute } from '@tanstack/react-router';
import { OrganizationsPage } from '@/pages/organizations/organizations-list';

export const Route = createFileRoute('/organizations')({
  component: OrganizationsPage,
});
