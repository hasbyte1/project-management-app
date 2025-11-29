import { useQuery } from '@tanstack/react-query';
import { useParams, useNavigate } from '@tanstack/react-router';
import { tasksApi } from '@/api/tasks';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { CreateTaskDialog } from '@/components/tasks/create-task-dialog';
import { Plus } from 'lucide-react';
import { getPriorityColor, getStatusColor, getInitials } from '@/lib/utils';

export function ListView() {
  const { organizationId, projectId } = useParams({ from: '/organizations/$organizationId/projects/$projectId/list' });
  const navigate = useNavigate();

  const { data: tasks, isLoading } = useQuery({
    queryKey: ['tasks', projectId],
    queryFn: () => tasksApi.list(projectId),
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-96">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Tasks</h2>
        <CreateTaskDialog projectId={projectId}>
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            New Task
          </Button>
        </CreateTaskDialog>
      </div>

      <div className="bg-white rounded-lg border">
        <div className="grid grid-cols-12 gap-4 px-4 py-3 border-b bg-muted/50 font-medium text-sm">
          <div className="col-span-5">Task</div>
          <div className="col-span-2">Status</div>
          <div className="col-span-2">Priority</div>
          <div className="col-span-2">Assignee</div>
          <div className="col-span-1">Due Date</div>
        </div>

        {tasks?.map((task) => (
          <div
            key={task.id}
            className="grid grid-cols-12 gap-4 px-4 py-3 border-b last:border-b-0 hover:bg-muted/50 cursor-pointer transition-colors"
            onClick={() => navigate({ to: `/organizations/${organizationId}/projects/${projectId}/tasks/${task.id}` })}
          >
            <div className="col-span-5">
              <div className="space-y-1">
                <h4 className="font-medium text-sm">{task.title}</h4>
                {task.description && (
                  <p className="text-xs text-muted-foreground line-clamp-1">
                    {task.description}
                  </p>
                )}
                {task.labels && task.labels.length > 0 && (
                  <div className="flex gap-1 flex-wrap">
                    {task.labels.map((label) => (
                      <Badge
                        key={label.id}
                        variant="outline"
                        className="text-xs"
                        style={{ borderColor: label.color, color: label.color }}
                      >
                        {label.name}
                      </Badge>
                    ))}
                  </div>
                )}
              </div>
            </div>

            <div className="col-span-2 flex items-center">
              {task.status && (
                <Badge className={`${getStatusColor(task.status.name)}`}>
                  {task.status.name}
                </Badge>
              )}
            </div>

            <div className="col-span-2 flex items-center">
              <Badge variant="outline" className={getPriorityColor(task.priority)}>
                {task.priority}
              </Badge>
            </div>

            <div className="col-span-2 flex items-center">
              {task.assignee ? (
                <div className="flex items-center gap-2">
                  <Avatar className="h-6 w-6">
                    <AvatarImage src={task.assignee.avatar_url} />
                    <AvatarFallback className="text-xs">
                      {getInitials(task.assignee.first_name, task.assignee.last_name)}
                    </AvatarFallback>
                  </Avatar>
                  <span className="text-sm">
                    {task.assignee.first_name} {task.assignee.last_name}
                  </span>
                </div>
              ) : (
                <span className="text-sm text-muted-foreground">Unassigned</span>
              )}
            </div>

            <div className="col-span-1 flex items-center">
              {task.due_date ? (
                <span className="text-sm">
                  {new Date(task.due_date).toLocaleDateString('en-US', {
                    month: 'short',
                    day: 'numeric',
                  })}
                </span>
              ) : (
                <span className="text-sm text-muted-foreground">-</span>
              )}
            </div>
          </div>
        ))}

        {tasks?.length === 0 && (
          <div className="text-center py-12">
            <p className="text-muted-foreground">No tasks found</p>
            <CreateTaskDialog projectId={projectId}>
              <Button className="mt-4">
                <Plus className="mr-2 h-4 w-4" />
                Create your first task
              </Button>
            </CreateTaskDialog>
          </div>
        )}
      </div>
    </div>
  );
}
