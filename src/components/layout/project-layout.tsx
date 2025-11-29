import { Outlet, useNavigate, useParams } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import { projectsApi } from '@/api/projects';
import { Button } from '@/components/ui/button';
import { LayoutList, LayoutGrid, Calendar, Settings } from 'lucide-react';

export function ProjectLayout() {
  const { organizationId, projectId } = useParams({ from: '/organizations/$organizationId/projects/$projectId/' });
  const navigate = useNavigate();

  const { data: project } = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => projectsApi.get(projectId),
  });

  return (
    <div className="min-h-screen">
      <header className="border-b bg-white">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between py-4">
            <div>
              <h1 className="text-2xl font-bold">{project?.name || 'Loading...'}</h1>
              <p className="text-sm text-muted-foreground">{project?.description}</p>
            </div>
            <Button variant="outline" size="icon">
              <Settings className="h-4 w-4" />
            </Button>
          </div>

          <nav className="flex gap-1">
            <Button
              variant="ghost"
              className="gap-2"
              onClick={() => navigate({ to: `/organizations/${organizationId}/projects/${projectId}/board` })}
            >
              <LayoutGrid className="h-4 w-4" />
              Board
            </Button>
            <Button
              variant="ghost"
              className="gap-2"
              onClick={() => navigate({ to: `/organizations/${organizationId}/projects/${projectId}/list` })}
            >
              <LayoutList className="h-4 w-4" />
              List
            </Button>
            <Button variant="ghost" className="gap-2">
              <Calendar className="h-4 w-4" />
              Calendar
            </Button>
          </nav>
        </div>
      </header>
      <main className="container mx-auto px-4 py-6">
        <Outlet />
      </main>
    </div>
  );
}
