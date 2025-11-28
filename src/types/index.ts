// Base types
export type UUID = string;

// User types
export interface User {
  id: UUID;
  email: string;
  first_name: string;
  last_name: string;
  avatar_url?: string;
  timezone: string;
  locale: string;
  email_verified: boolean;
  is_active: boolean;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface UserToken {
  id: UUID;
  user_id: UUID;
  token_hash: string;
  token_type: 'refresh' | 'reset_password' | 'email_verification';
  expires_at: string;
  created_at: string;
  used_at?: string;
}

// Organization types
export type OrgRole = 'owner' | 'admin' | 'member';

export interface Organization {
  id: UUID;
  parent_id?: UUID;
  name: string;
  slug: string;
  description?: string;
  logo_url?: string;
  depth: number;
  path?: string;
  settings: Record<string, any>;
  created_by: UUID;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface OrganizationMember {
  id: UUID;
  organization_id: UUID;
  user_id: UUID;
  role: OrgRole;
  invited_by?: UUID;
  invited_at?: string;
  joined_at: string;
  created_at: string;
  updated_at: string;
  user?: User;
}

// Team types
export interface Team {
  id: UUID;
  organization_id: UUID;
  name: string;
  description?: string;
  color?: string;
  created_by: UUID;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface TeamMember {
  id: UUID;
  team_id: UUID;
  user_id: UUID;
  added_by?: UUID;
  created_at: string;
  user?: User;
}

// Project types
export type ProjectVisibility = 'private' | 'team' | 'organization';
export type ProjectStatus = 'active' | 'on_hold' | 'archived' | 'completed';
export type ProjectRole = 'owner' | 'editor' | 'viewer';

export interface Project {
  id: UUID;
  organization_id: UUID;
  team_id?: UUID;
  name: string;
  description?: string;
  key?: string;
  color?: string;
  icon?: string;
  visibility: ProjectVisibility;
  status: ProjectStatus;
  start_date?: string;
  due_date?: string;
  settings: Record<string, any>;
  created_by: UUID;
  created_at: string;
  updated_at: string;
  archived_at?: string;
  deleted_at?: string;
}

export interface ProjectMember {
  id: UUID;
  project_id: UUID;
  user_id: UUID;
  role: ProjectRole;
  added_by?: UUID;
  created_at: string;
  updated_at: string;
  user?: User;
}

// Task types
export type TaskPriority = 'urgent' | 'high' | 'medium' | 'low' | 'none';
export type DependencyType = 'blocks' | 'blocked_by' | 'relates_to' | 'duplicates';

export interface TaskStatus {
  id: UUID;
  project_id: UUID;
  name: string;
  color: string;
  position: number;
  is_default: boolean;
  is_completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: UUID;
  project_id: UUID;
  parent_task_id?: UUID;
  title: string;
  description?: string;
  task_number: number;
  status_id: UUID;
  priority: TaskPriority;
  assignee_id?: UUID;
  reporter_id: UUID;
  start_date?: string;
  due_date?: string;
  completed_at?: string;
  estimated_hours?: number;
  actual_hours: number;
  position: number;
  custom_fields: Record<string, any>;
  created_by: UUID;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
  // Relations
  status?: TaskStatus;
  assignee?: User;
  reporter?: User;
  labels?: Label[];
  subtasks?: Task[];
}

export interface TaskAssignee {
  id: UUID;
  task_id: UUID;
  user_id: UUID;
  assigned_by?: UUID;
  created_at: string;
  user?: User;
}

export interface TaskDependency {
  id: UUID;
  task_id: UUID;
  depends_on_task_id: UUID;
  dependency_type: DependencyType;
  created_by: UUID;
  created_at: string;
}

// Label types
export interface Label {
  id: UUID;
  organization_id: UUID;
  project_id?: UUID;
  name: string;
  color: string;
  description?: string;
  created_by: UUID;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface TaskLabel {
  id: UUID;
  task_id: UUID;
  label_id: UUID;
  created_at: string;
  label?: Label;
}

// Comment types
export interface Comment {
  id: UUID;
  task_id: UUID;
  user_id: UUID;
  content: string;
  parent_comment_id?: UUID;
  is_edited: boolean;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
  user?: User;
  replies?: Comment[];
}

// Activity log types
export type ActivityType =
  | 'task_created'
  | 'task_updated'
  | 'task_deleted'
  | 'task_completed'
  | 'task_assigned'
  | 'task_unassigned'
  | 'task_status_changed'
  | 'comment_added'
  | 'comment_updated'
  | 'comment_deleted'
  | 'attachment_added'
  | 'attachment_removed'
  | 'dependency_added'
  | 'dependency_removed'
  | 'label_added'
  | 'label_removed';

export interface ActivityLog {
  id: UUID;
  organization_id: UUID;
  project_id?: UUID;
  task_id?: UUID;
  user_id: UUID;
  activity_type: ActivityType;
  metadata: Record<string, any>;
  created_at: string;
  user?: User;
}

// Attachment types
export interface Attachment {
  id: UUID;
  task_id?: UUID;
  comment_id?: UUID;
  uploaded_by: UUID;
  file_name: string;
  file_size: number;
  mime_type: string;
  storage_path: string;
  thumbnail_path?: string;
  created_at: string;
  deleted_at?: string;
  uploader?: User;
}

// Time tracking types
export interface TimeEntry {
  id: UUID;
  task_id: UUID;
  user_id: UUID;
  description?: string;
  hours: number;
  started_at: string;
  ended_at?: string;
  is_billable: boolean;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
  user?: User;
}

// Custom field types
export type CustomFieldType =
  | 'text'
  | 'number'
  | 'date'
  | 'dropdown'
  | 'multi_select'
  | 'checkbox'
  | 'url'
  | 'email'
  | 'user'
  | 'currency';

export interface CustomField {
  id: UUID;
  project_id: UUID;
  name: string;
  field_type: CustomFieldType;
  description?: string;
  is_required: boolean;
  options?: any[];
  default_value?: any;
  position: number;
  created_by: UUID;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

// View types
export type ViewType = 'list' | 'board' | 'calendar' | 'timeline' | 'table';

export interface View {
  id: UUID;
  project_id: UUID;
  created_by: UUID;
  name: string;
  view_type: ViewType;
  is_default: boolean;
  is_shared: boolean;
  filters: Record<string, any>;
  sort_by: any[];
  group_by?: string;
  visible_fields: string[];
  settings: Record<string, any>;
  position: number;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

// Notification types
export type NotificationType =
  | 'task_assigned'
  | 'task_mentioned'
  | 'task_due_soon'
  | 'task_overdue'
  | 'comment_added'
  | 'comment_mentioned'
  | 'status_changed'
  | 'dependency_blocked';

export interface Notification {
  id: UUID;
  user_id: UUID;
  notification_type: NotificationType;
  title: string;
  message?: string;
  task_id?: UUID;
  project_id?: UUID;
  triggered_by?: UUID;
  is_read: boolean;
  read_at?: string;
  created_at: string;
  task?: Task;
  project?: Project;
  triggerer?: User;
}

export interface NotificationPreference {
  id: UUID;
  user_id: UUID;
  notification_type: NotificationType;
  email_enabled: boolean;
  push_enabled: boolean;
  in_app_enabled: boolean;
}

// Template types
export interface ProjectTemplate {
  id: UUID;
  organization_id: UUID;
  name: string;
  description?: string;
  is_public: boolean;
  template_data: Record<string, any>;
  created_by: UUID;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface TaskTemplate {
  id: UUID;
  project_id?: UUID;
  organization_id?: UUID;
  name: string;
  description?: string;
  template_data: Record<string, any>;
  created_by: UUID;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

// Automation types
export interface AutomationRule {
  id: UUID;
  project_id: UUID;
  name: string;
  description?: string;
  is_active: boolean;
  trigger_event: string;
  conditions: Record<string, any>;
  actions: Record<string, any>;
  created_by: UUID;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

// Auth types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}
