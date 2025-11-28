import { createFileRoute } from '@tanstack/react-router';
import { ListView } from '@/pages/tasks/list-view';

export const Route = createFileRoute('/organizations/$organizationId/projects/$projectId/list')({
  component: ListView,
});
