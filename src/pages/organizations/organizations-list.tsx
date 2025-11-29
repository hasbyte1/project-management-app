import { useQuery } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';
import { organizationsApi } from '@/api/organizations';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { CreateOrganizationDialog } from '@/components/organizations/create-organization-dialog';
import { Building2, Plus } from 'lucide-react';

export function OrganizationsPage() {
  const navigate = useNavigate();

  const { data: organizations, isLoading } = useQuery({
    queryKey: ['organizations'],
    queryFn: () => organizationsApi.list(),
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8 px-4">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold">Organizations</h1>
          <p className="text-muted-foreground mt-2">
            Select an organization to get started
          </p>
        </div>
        <CreateOrganizationDialog>
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            New Organization
          </Button>
        </CreateOrganizationDialog>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {organizations?.map((org) => (
          <Card
            key={org.id}
            className="cursor-pointer hover:shadow-lg transition-shadow"
            onClick={() => navigate({ to: `/organizations/${org.id}/projects` })}
          >
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="p-2 bg-primary/10 rounded-lg">
                  <Building2 className="h-6 w-6 text-primary" />
                </div>
                <div className="flex-1">
                  <CardTitle>{org.name}</CardTitle>
                  <CardDescription>{org.slug}</CardDescription>
                </div>
              </div>
            </CardHeader>
            {org.description && (
              <CardContent>
                <p className="text-sm text-muted-foreground">{org.description}</p>
              </CardContent>
            )}
          </Card>
        ))}
      </div>

      {organizations?.length === 0 && (
        <div className="text-center py-12">
          <Building2 className="mx-auto h-12 w-12 text-muted-foreground" />
          <h3 className="mt-4 text-lg font-semibold">No organizations found</h3>
          <p className="text-muted-foreground mt-2">
            Create your first organization to get started
          </p>
          <CreateOrganizationDialog>
            <Button className="mt-4">
              <Plus className="mr-2 h-4 w-4" />
              Create Organization
            </Button>
          </CreateOrganizationDialog>
        </div>
      )}
    </div>
  );
}
