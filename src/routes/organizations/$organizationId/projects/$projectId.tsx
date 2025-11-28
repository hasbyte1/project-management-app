import { createFileRoute } from '@tanstack/react-router';
import { ProjectLayout } from '@/components/layout/project-layout';

export const Route = createFileRoute('/organizations/$organizationId/projects/$projectId')({
  component: ProjectLayout,
});
