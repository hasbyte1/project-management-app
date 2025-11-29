import { ProjectLayout } from '@/components/layout/project-layout'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute(
  '/organizations/$organizationId/projects/$projectId/',
)({
  component: ProjectLayout,
})
