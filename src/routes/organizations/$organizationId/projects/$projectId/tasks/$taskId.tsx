import { createFileRoute } from '@tanstack/react-router';
import { TaskDetailPage } from '@/pages/tasks/task-detail';

export const Route = createFileRoute('/organizations/$organizationId/projects/$projectId/tasks/$taskId')({
  component: TaskDetailPage,
});
