import { Outlet } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import { useParams } from '@tanstack/react-router';
import { organizationsApi } from '@/api/organizations';

export function OrganizationLayout() {
  const { organizationId } = useParams({ from: '/organizations/$organizationId' });

  const { data: organization } = useQuery({
    queryKey: ['organization', organizationId],
    queryFn: () => organizationsApi.get(organizationId),
  });

  return (
    <div className="min-h-screen">
      <header className="border-b bg-white">
        <div className="container mx-auto px-4 py-4">
          <h2 className="text-xl font-semibold">{organization?.name || 'Loading...'}</h2>
        </div>
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  );
}
