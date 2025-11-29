import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useParams, useNavigate } from '@tanstack/react-router';
import { tasksApi } from '@/api/tasks';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  ArrowLeft,
  Calendar,
  User,
  Flag,
  MessageSquare,
  Paperclip,
  Clock,
} from 'lucide-react';
import { getInitials, formatDate, formatDateTime } from '@/lib/utils';
import { useState } from 'react';
import { toast } from 'sonner';

export function TaskDetailPage() {
  const { organizationId, projectId, taskId } = useParams({
    from: '/organizations/$organizationId/projects/$projectId/tasks/$taskId',
  });
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [newComment, setNewComment] = useState('');

  const { data: task, isLoading } = useQuery({
    queryKey: ['task', taskId],
    queryFn: () => tasksApi.get(taskId),
  });

  const { data: comments } = useQuery({
    queryKey: ['comments', taskId],
    queryFn: () => tasksApi.getComments(taskId),
  });

  const { data: timeEntries } = useQuery({
    queryKey: ['timeEntries', taskId],
    queryFn: () => tasksApi.getTimeEntries(taskId),
  });

  const { data: attachments } = useQuery({
    queryKey: ['attachments', taskId],
    queryFn: () => tasksApi.getAttachments(taskId),
  });

  const { data: statuses } = useQuery({
    queryKey: ['statuses', projectId],
    queryFn: () => tasksApi.getStatuses(projectId),
  });

  const updateTaskMutation = useMutation({
    mutationFn: (data: any) => tasksApi.update(taskId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['task', taskId] });
      queryClient.invalidateQueries({ queryKey: ['tasks', projectId] });
      toast.success('Task updated successfully');
    },
  });

  const createCommentMutation = useMutation({
    mutationFn: (content: string) => tasksApi.createComment(taskId, content),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', taskId] });
      setNewComment('');
      toast.success('Comment added');
    },
  });

  const handleAddComment = () => {
    if (newComment.trim()) {
      createCommentMutation.mutate(newComment);
    }
  };

  if (isLoading || !task) {
    return (
      <div className="flex items-center justify-center min-h-96">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    );
  }

  const totalTimeLogged = timeEntries?.reduce((sum, entry) => sum + Number(entry.hours), 0) || 0;

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-6">
        <Button
          variant="ghost"
          onClick={() => navigate({ to: `/organizations/${organizationId}/projects/${projectId}/board` })}
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Board
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Task Header */}
          <div>
            <h1 className="text-3xl font-bold mb-2">{task.title}</h1>
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <span>Task #{task.task_number}</span>
              <span>•</span>
              <span>Created {formatDate(task.created_at)}</span>
            </div>
          </div>

          {/* Description */}
          <Card>
            <CardHeader>
              <CardTitle>Description</CardTitle>
            </CardHeader>
            <CardContent>
              {task.description ? (
                <p className="text-sm whitespace-pre-wrap">{task.description}</p>
              ) : (
                <p className="text-sm text-muted-foreground">No description provided</p>
              )}
            </CardContent>
          </Card>

          {/* Comments */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <MessageSquare className="h-5 w-5" />
                Comments ({comments?.length || 0})
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {comments?.map((comment) => (
                <div key={comment.id} className="space-y-2">
                  <div className="flex items-start gap-3">
                    <Avatar className="h-8 w-8">
                      <AvatarImage src={comment.user?.avatar_url} />
                      <AvatarFallback>
                        {comment.user &&
                          getInitials(comment.user.first_name, comment.user.last_name)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <span className="font-medium text-sm">
                          {comment.user?.first_name} {comment.user?.last_name}
                        </span>
                        <span className="text-xs text-muted-foreground">
                          {formatDateTime(comment.created_at)}
                        </span>
                        {comment.is_edited && (
                          <Badge variant="outline" className="text-xs">
                            Edited
                          </Badge>
                        )}
                      </div>
                      <p className="text-sm mt-1 whitespace-pre-wrap">{comment.content}</p>
                    </div>
                  </div>
                  <Separator />
                </div>
              ))}

              {/* Add Comment */}
              <div className="space-y-3">
                <Textarea
                  placeholder="Add a comment..."
                  value={newComment}
                  onChange={(e) => setNewComment(e.target.value)}
                  rows={3}
                />
                <Button
                  onClick={handleAddComment}
                  disabled={!newComment.trim() || createCommentMutation.isPending}
                >
                  {createCommentMutation.isPending ? 'Adding...' : 'Add Comment'}
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Attachments */}
          {attachments && attachments.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Paperclip className="h-5 w-5" />
                  Attachments ({attachments.length})
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {attachments.map((attachment) => (
                    <div
                      key={attachment.id}
                      className="flex items-center justify-between p-3 border rounded-lg hover:bg-muted/50"
                    >
                      <div className="flex items-center gap-3">
                        <Paperclip className="h-4 w-4 text-muted-foreground" />
                        <div>
                          <p className="text-sm font-medium">{attachment.file_name}</p>
                          <p className="text-xs text-muted-foreground">
                            {(attachment.file_size / 1024).toFixed(2)} KB
                          </p>
                        </div>
                      </div>
                      <Button variant="outline" size="sm">
                        Download
                      </Button>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Status */}
          <Card>
            <CardContent className="pt-6 space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium flex items-center gap-2">
                  <Flag className="h-4 w-4" />
                  Status
                </label>
                <Select
                  value={task.status_id}
                  onValueChange={(value) =>
                    updateTaskMutation.mutate({ status_id: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {statuses?.map((status) => (
                      <SelectItem key={status.id} value={status.id}>
                        <div className="flex items-center gap-2">
                          <div
                            className="w-2 h-2 rounded-full"
                            style={{ backgroundColor: status.color }}
                          />
                          {status.name}
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Priority</label>
                <Select
                  value={task.priority}
                  onValueChange={(value) =>
                    updateTaskMutation.mutate({ priority: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="urgent">Urgent</SelectItem>
                    <SelectItem value="high">High</SelectItem>
                    <SelectItem value="medium">Medium</SelectItem>
                    <SelectItem value="low">Low</SelectItem>
                    <SelectItem value="none">None</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <Separator />

              <div className="space-y-2">
                <label className="text-sm font-medium flex items-center gap-2">
                  <User className="h-4 w-4" />
                  Assignee
                </label>
                {task.assignee ? (
                  <div className="flex items-center gap-2">
                    <Avatar className="h-8 w-8">
                      <AvatarImage src={task.assignee.avatar_url} />
                      <AvatarFallback>
                        {getInitials(task.assignee.first_name, task.assignee.last_name)}
                      </AvatarFallback>
                    </Avatar>
                    <span className="text-sm">
                      {task.assignee.first_name} {task.assignee.last_name}
                    </span>
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground">Unassigned</p>
                )}
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium flex items-center gap-2">
                  <Calendar className="h-4 w-4" />
                  Due Date
                </label>
                <p className="text-sm">
                  {task.due_date ? formatDate(task.due_date) : 'No due date'}
                </p>
              </div>

              <Separator />

              <div className="space-y-2">
                <label className="text-sm font-medium flex items-center gap-2">
                  <Clock className="h-4 w-4" />
                  Time Tracking
                </label>
                <div className="space-y-1 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Estimated:</span>
                    <span>{task.estimated_hours || 0}h</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Logged:</span>
                    <span>{totalTimeLogged.toFixed(2)}h</span>
                  </div>
                  {task.estimated_hours && (
                    <div className="flex justify-between font-medium">
                      <span className="text-muted-foreground">Remaining:</span>
                      <span>
                        {Math.max(0, task.estimated_hours - totalTimeLogged).toFixed(2)}h
                      </span>
                    </div>
                  )}
                </div>
              </div>

              {task.labels && task.labels.length > 0 && (
                <>
                  <Separator />
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Labels</label>
                    <div className="flex flex-wrap gap-2">
                      {task.labels.map((label) => (
                        <Badge
                          key={label.id}
                          variant="outline"
                          style={{ borderColor: label.color, color: label.color }}
                        >
                          {label.name}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
