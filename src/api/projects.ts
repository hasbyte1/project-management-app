import { apiClient } from '@/lib/api-client';
import type { Project, ProjectMember, UUID } from '@/types';

export interface CreateProjectRequest {
  organization_id: UUID;
  team_id?: UUID;
  name: string;
  description?: string;
  key?: string;
  color?: string;
  icon?: string;
  visibility?: 'private' | 'team' | 'organization';
  start_date?: string;
  due_date?: string;
}

export interface UpdateProjectRequest extends Partial<CreateProjectRequest> {
  status?: 'active' | 'on_hold' | 'archived' | 'completed';
}

export const projectsApi = {
  list: async (organizationId: UUID): Promise<Project[]> => {
    const response = await apiClient.get<Project[]>(`/organizations/${organizationId}/projects`);
    return response.data;
  },

  get: async (projectId: UUID): Promise<Project> => {
    const response = await apiClient.get<Project>(`/projects/${projectId}`);
    return response.data;
  },

  create: async (data: CreateProjectRequest): Promise<Project> => {
    const response = await apiClient.post<Project>('/projects', data);
    return response.data;
  },

  update: async (projectId: UUID, data: UpdateProjectRequest): Promise<Project> => {
    const response = await apiClient.patch<Project>(`/projects/${projectId}`, data);
    return response.data;
  },

  delete: async (projectId: UUID): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}`);
  },

  archive: async (projectId: UUID): Promise<Project> => {
    const response = await apiClient.post<Project>(`/projects/${projectId}/archive`);
    return response.data;
  },

  unarchive: async (projectId: UUID): Promise<Project> => {
    const response = await apiClient.post<Project>(`/projects/${projectId}/unarchive`);
    return response.data;
  },

  // Project members
  getMembers: async (projectId: UUID): Promise<ProjectMember[]> => {
    const response = await apiClient.get<ProjectMember[]>(`/projects/${projectId}/members`);
    return response.data;
  },

  addMember: async (projectId: UUID, userId: UUID, role: 'owner' | 'editor' | 'viewer'): Promise<ProjectMember> => {
    const response = await apiClient.post<ProjectMember>(`/projects/${projectId}/members`, {
      user_id: userId,
      role,
    });
    return response.data;
  },

  updateMemberRole: async (projectId: UUID, memberId: UUID, role: 'owner' | 'editor' | 'viewer'): Promise<ProjectMember> => {
    const response = await apiClient.patch<ProjectMember>(`/projects/${projectId}/members/${memberId}`, {
      role,
    });
    return response.data;
  },

  removeMember: async (projectId: UUID, memberId: UUID): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/members/${memberId}`);
  },
};
