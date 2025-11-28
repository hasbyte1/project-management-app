import { createFileRoute } from '@tanstack/react-router';
import { BoardView } from '@/pages/tasks/board-view';

export const Route = createFileRoute('/organizations/$organizationId/projects/$projectId/board')({
  component: BoardView,
});
