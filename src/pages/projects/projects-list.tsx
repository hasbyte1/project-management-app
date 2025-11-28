import { useQuery } from '@tanstack/react-query';
import { useNavigate, useParams } from '@tanstack/react-router';
import { projectsApi } from '@/api/projects';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { FolderKanban, Plus } from 'lucide-react';

export function ProjectsPage() {
  const { organizationId } = useParams({ from: '/organizations/$organizationId/projects' });
  const navigate = useNavigate();

  const { data: projects, isLoading } = useQuery({
    queryKey: ['projects', organizationId],
    queryFn: () => projectsApi.list(organizationId),
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
          <h1 className="text-3xl font-bold">Projects</h1>
          <p className="text-muted-foreground mt-2">
            Manage your projects and track progress
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          New Project
        </Button>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {projects?.map((project) => (
          <Card
            key={project.id}
            className="cursor-pointer hover:shadow-lg transition-shadow"
            onClick={() => navigate({ to: `/organizations/${organizationId}/projects/${project.id}/board` })}
          >
            <CardHeader>
              <div className="flex items-center gap-3">
                <div
                  className="p-2 rounded-lg"
                  style={{ backgroundColor: project.color || '#3b82f6' + '20' }}
                >
                  <FolderKanban
                    className="h-6 w-6"
                    style={{ color: project.color || '#3b82f6' }}
                  />
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <CardTitle>{project.name}</CardTitle>
                    {project.key && (
                      <Badge variant="outline" className="text-xs">
                        {project.key}
                      </Badge>
                    )}
                  </div>
                  <div className="flex gap-2 mt-1">
                    <Badge variant="secondary" className="text-xs capitalize">
                      {project.status}
                    </Badge>
                    <Badge variant="outline" className="text-xs capitalize">
                      {project.visibility}
                    </Badge>
                  </div>
                </div>
              </div>
            </CardHeader>
            {project.description && (
              <CardContent>
                <p className="text-sm text-muted-foreground line-clamp-2">
                  {project.description}
                </p>
              </CardContent>
            )}
          </Card>
        ))}
      </div>

      {projects?.length === 0 && (
        <div className="text-center py-12">
          <FolderKanban className="mx-auto h-12 w-12 text-muted-foreground" />
          <h3 className="mt-4 text-lg font-semibold">No projects found</h3>
          <p className="text-muted-foreground mt-2">
            Create your first project to start managing tasks
          </p>
          <Button className="mt-4">
            <Plus className="mr-2 h-4 w-4" />
            Create Project
          </Button>
        </div>
      )}
    </div>
  );
}
