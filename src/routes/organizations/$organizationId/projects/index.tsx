import { ProjectsPage } from '@/pages/projects/projects-list'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute(
  '/organizations/$organizationId/projects/',
)({
  component: ProjectsPage,
})
