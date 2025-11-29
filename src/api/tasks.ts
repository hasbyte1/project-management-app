import { apiClient } from '@/lib/api-client';
import type { Task, TaskStatus, Comment, TimeEntry, Attachment, UUID, TaskPriority } from '@/types';

export interface CreateTaskRequest {
  project_id: UUID;
  parent_task_id?: UUID;
  title: string;
  description?: string;
  status_id: UUID;
  priority?: TaskPriority;
  assignee_id?: UUID;
  start_date?: string;
  due_date?: string;
  estimated_hours?: number;
  custom_fields?: Record<string, any>;
}

export interface UpdateTaskRequest extends Partial<CreateTaskRequest> {
  position?: number;
}

export interface TaskFilters {
  status_id?: UUID[];
  assignee_id?: UUID[];
  priority?: TaskPriority[];
  labels?: UUID[];
  search?: string;
  due_date_from?: string;
  due_date_to?: string;
}

export const tasksApi = {
  // Tasks
  list: async (projectId: UUID, filters?: TaskFilters): Promise<Task[]> => {
    const response = await apiClient.get<Task[]>(`/projects/${projectId}/tasks`, {
      params: filters,
    });
    return response.data;
  },

  get: async (taskId: UUID): Promise<Task> => {
    const response = await apiClient.get<Task>(`/tasks/${taskId}`);
    return response.data;
  },

  create: async (data: CreateTaskRequest): Promise<Task> => {
    const response = await apiClient.post<Task>('/tasks', data);
    return response.data;
  },

  update: async (taskId: UUID, data: UpdateTaskRequest): Promise<Task> => {
    const response = await apiClient.patch<Task>(`/tasks/${taskId}`, data);
    return response.data;
  },

  delete: async (taskId: UUID): Promise<void> => {
    await apiClient.delete(`/tasks/${taskId}`);
  },

  updateStatus: async (taskId: UUID, statusId: UUID): Promise<Task> => {
    const response = await apiClient.patch<Task>(`/tasks/${taskId}/status`, {
      status_id: statusId,
    });
    return response.data;
  },

  updatePosition: async (taskId: UUID, position: number): Promise<Task> => {
    const response = await apiClient.patch<Task>(`/tasks/${taskId}/position`, {
      position,
    });
    return response.data;
  },

  // Task statuses
  getStatuses: async (projectId: UUID): Promise<TaskStatus[]> => {
    const response = await apiClient.get<TaskStatus[]>(`/projects/${projectId}/statuses`);
    return response.data;
  },

  createStatus: async (projectId: UUID, name: string, color: string, position: number): Promise<TaskStatus> => {
    const response = await apiClient.post<TaskStatus>(`/projects/${projectId}/statuses`, {
      name,
      color,
      position,
    });
    return response.data;
  },

  updateTaskStatus: async (statusId: UUID, data: Partial<TaskStatus>): Promise<TaskStatus> => {
    const response = await apiClient.patch<TaskStatus>(`/statuses/${statusId}`, data);
    return response.data;
  },

  deleteStatus: async (statusId: UUID): Promise<void> => {
    await apiClient.delete(`/statuses/${statusId}`);
  },

  // Comments
  getComments: async (taskId: UUID): Promise<Comment[]> => {
    const response = await apiClient.get<Comment[]>(`/tasks/${taskId}/comments`);
    return response.data;
  },

  createComment: async (taskId: UUID, content: string, parentCommentId?: UUID): Promise<Comment> => {
    const response = await apiClient.post<Comment>(`/tasks/${taskId}/comments`, {
      content,
      parent_comment_id: parentCommentId,
    });
    return response.data;
  },

  updateComment: async (commentId: UUID, content: string): Promise<Comment> => {
    const response = await apiClient.patch<Comment>(`/comments/${commentId}`, {
      content,
    });
    return response.data;
  },

  deleteComment: async (commentId: UUID): Promise<void> => {
    await apiClient.delete(`/comments/${commentId}`);
  },

  // Time entries
  getTimeEntries: async (taskId: UUID): Promise<TimeEntry[]> => {
    const response = await apiClient.get<TimeEntry[]>(`/tasks/${taskId}/time-entries`);
    return response.data;
  },

  createTimeEntry: async (
    taskId: UUID,
    data: {
      hours: number;
      description?: string;
      started_at: string;
      ended_at?: string;
      is_billable?: boolean;
    }
  ): Promise<TimeEntry> => {
    const response = await apiClient.post<TimeEntry>(`/tasks/${taskId}/time-entries`, data);
    return response.data;
  },

  updateTimeEntry: async (timeEntryId: UUID, data: Partial<TimeEntry>): Promise<TimeEntry> => {
    const response = await apiClient.patch<TimeEntry>(`/time-entries/${timeEntryId}`, data);
    return response.data;
  },

  deleteTimeEntry: async (timeEntryId: UUID): Promise<void> => {
    await apiClient.delete(`/time-entries/${timeEntryId}`);
  },

  // Attachments
  getAttachments: async (taskId: UUID): Promise<Attachment[]> => {
    const response = await apiClient.get<Attachment[]>(`/tasks/${taskId}/attachments`);
    return response.data;
  },

  uploadAttachment: async (taskId: UUID, file: File): Promise<Attachment> => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await apiClient.post<Attachment>(`/tasks/${taskId}/attachments`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  deleteAttachment: async (attachmentId: UUID): Promise<void> => {
    await apiClient.delete(`/attachments/${attachmentId}`);
  },

  // Labels
  addLabel: async (taskId: UUID, labelId: UUID): Promise<void> => {
    await apiClient.post(`/tasks/${taskId}/labels`, { label_id: labelId });
  },

  removeLabel: async (taskId: UUID, labelId: UUID): Promise<void> => {
    await apiClient.delete(`/tasks/${taskId}/labels/${labelId}`);
  },

  // Dependencies
  addDependency: async (taskId: UUID, dependsOnTaskId: UUID, dependencyType: string): Promise<void> => {
    await apiClient.post(`/tasks/${taskId}/dependencies`, {
      depends_on_task_id: dependsOnTaskId,
      dependency_type: dependencyType,
    });
  },

  removeDependency: async (taskId: UUID, dependencyId: UUID): Promise<void> => {
    await apiClient.delete(`/tasks/${taskId}/dependencies/${dependencyId}`);
  },
};
