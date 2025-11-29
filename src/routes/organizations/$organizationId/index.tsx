import { OrganizationLayout } from '@/components/layout/organization-layout'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/organizations/$organizationId/')({
  component: OrganizationLayout,
})
