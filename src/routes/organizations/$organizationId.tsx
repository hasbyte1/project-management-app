import { createFileRoute } from '@tanstack/react-router';
import { OrganizationLayout } from '@/components/layout/organization-layout';

export const Route = createFileRoute('/organizations/$organizationId')({
  component: OrganizationLayout,
});
