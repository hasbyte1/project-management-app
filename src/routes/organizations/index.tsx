import { OrganizationsPage } from '@/pages/organizations/organizations-list'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/organizations/')({
  component: OrganizationsPage,
})
