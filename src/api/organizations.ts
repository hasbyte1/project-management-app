import { apiClient } from '@/lib/api-client';
import type { Organization, OrganizationMember, Team, TeamMember, Label, UUID } from '@/types';

export interface CreateOrganizationRequest {
  name: string;
  slug: string;
  description?: string;
  logo_url?: string;
  parent_id?: UUID;
}

export const organizationsApi = {
  // Organizations
  list: async (): Promise<Organization[]> => {
    const response = await apiClient.get<Organization[]>('/organizations');
    return response.data;
  },

  get: async (organizationId: UUID): Promise<Organization> => {
    const response = await apiClient.get<Organization>(`/organizations/${organizationId}`);
    return response.data;
  },

  create: async (data: CreateOrganizationRequest): Promise<Organization> => {
    const response = await apiClient.post<Organization>('/organizations', data);
    return response.data;
  },

  update: async (organizationId: UUID, data: Partial<CreateOrganizationRequest>): Promise<Organization> => {
    const response = await apiClient.patch<Organization>(`/organizations/${organizationId}`, data);
    return response.data;
  },

  delete: async (organizationId: UUID): Promise<void> => {
    await apiClient.delete(`/organizations/${organizationId}`);
  },

  // Organization members
  getMembers: async (organizationId: UUID): Promise<OrganizationMember[]> => {
    const response = await apiClient.get<OrganizationMember[]>(`/organizations/${organizationId}/members`);
    return response.data;
  },

  addMember: async (organizationId: UUID, email: string, role: 'owner' | 'admin' | 'member'): Promise<OrganizationMember> => {
    const response = await apiClient.post<OrganizationMember>(`/organizations/${organizationId}/members`, {
      email,
      role,
    });
    return response.data;
  },

  updateMemberRole: async (organizationId: UUID, memberId: UUID, role: 'owner' | 'admin' | 'member'): Promise<OrganizationMember> => {
    const response = await apiClient.patch<OrganizationMember>(`/organizations/${organizationId}/members/${memberId}`, {
      role,
    });
    return response.data;
  },

  removeMember: async (organizationId: UUID, memberId: UUID): Promise<void> => {
    await apiClient.delete(`/organizations/${organizationId}/members/${memberId}`);
  },

  // Teams
  getTeams: async (organizationId: UUID): Promise<Team[]> => {
    const response = await apiClient.get<Team[]>(`/organizations/${organizationId}/teams`);
    return response.data;
  },

  createTeam: async (organizationId: UUID, name: string, description?: string, color?: string): Promise<Team> => {
    const response = await apiClient.post<Team>(`/organizations/${organizationId}/teams`, {
      name,
      description,
      color,
    });
    return response.data;
  },

  updateTeam: async (teamId: UUID, data: Partial<Team>): Promise<Team> => {
    const response = await apiClient.patch<Team>(`/teams/${teamId}`, data);
    return response.data;
  },

  deleteTeam: async (teamId: UUID): Promise<void> => {
    await apiClient.delete(`/teams/${teamId}`);
  },

  // Team members
  getTeamMembers: async (teamId: UUID): Promise<TeamMember[]> => {
    const response = await apiClient.get<TeamMember[]>(`/teams/${teamId}/members`);
    return response.data;
  },

  addTeamMember: async (teamId: UUID, userId: UUID): Promise<TeamMember> => {
    const response = await apiClient.post<TeamMember>(`/teams/${teamId}/members`, {
      user_id: userId,
    });
    return response.data;
  },

  removeTeamMember: async (teamId: UUID, memberId: UUID): Promise<void> => {
    await apiClient.delete(`/teams/${teamId}/members/${memberId}`);
  },

  // Labels
  getLabels: async (organizationId: UUID, projectId?: UUID): Promise<Label[]> => {
    const response = await apiClient.get<Label[]>(`/organizations/${organizationId}/labels`, {
      params: projectId ? { project_id: projectId } : undefined,
    });
    return response.data;
  },

  createLabel: async (organizationId: UUID, data: { name: string; color: string; description?: string; project_id?: UUID }): Promise<Label> => {
    const response = await apiClient.post<Label>(`/organizations/${organizationId}/labels`, data);
    return response.data;
  },

  updateLabel: async (labelId: UUID, data: Partial<Label>): Promise<Label> => {
    const response = await apiClient.patch<Label>(`/labels/${labelId}`, data);
    return response.data;
  },

  deleteLabel: async (labelId: UUID): Promise<void> => {
    await apiClient.delete(`/labels/${labelId}`);
  },
};
