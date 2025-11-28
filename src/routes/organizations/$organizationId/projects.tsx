import { createFileRoute } from '@tanstack/react-router';
import { ProjectsPage } from '@/pages/projects/projects-list';

export const Route = createFileRoute('/organizations/$organizationId/projects')({
  component: ProjectsPage,
});
