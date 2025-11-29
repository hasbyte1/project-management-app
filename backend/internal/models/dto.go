package models

import (
	"time"

	"github.com/google/uuid"
)

// Auth DTOs
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Organization DTOs

type OrganizationDTO struct {
	ID          string  `json:"id"`
	ParentID    *string `json:"parent_id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	LogoURL     *string `json:"logo_url"`
	Depth       int     `json:"depth"`
	Path        *string `json:"path"`
	Settings    []byte  `json:"settings"`
	CreatedBy   string  `json:"created_by"`
}
type CreateOrganizationRequest struct {
	Name        string     `json:"name" validate:"required"`
	Slug        string     `json:"slug" validate:"required,alphanum"`
	Description *string    `json:"description,omitempty"`
	LogoURL     *string    `json:"logo_url,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
}

type UpdateOrganizationRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	LogoURL     *string `json:"logo_url,omitempty"`
}

type AddOrganizationMemberRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=owner admin member"`
}

// Project DTOs

type ProjectFilters struct {
	OrganizationIDs []uuid.UUID `json:"organization_ids,omitempty"`
	TeamIDs         []uuid.UUID `json:"team_ids,omitempty"`
	Key             string      `json:"assignee_id,omitempty"`
	Search          string      `json:"search,omitempty"`
	Visibility      string      `json:"visibility,omitempty"`
	StartDate       *time.Time  `json:"start_date,omitempty"`
	DueDate         *time.Time  `json:"due_date,omitempty"`
}

type ProjectDTO struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	TeamID         *string    `json:"team_id"`
	Name           string     `json:"name"`
	Description    *string    `json:"description"`
	Key            *string    `json:"key"`
	Color          *string    `json:"color"`
	Icon           *string    `json:"icon"`
	Visibility     string     `json:"visibility"`
	StartDate      *time.Time `json:"start_date"`
	DueDate        *time.Time `json:"due_date"`
}

type CreateProjectRequest struct {
	OrganizationID uuid.UUID  `json:"organization_id" validate:"required"`
	TeamID         *uuid.UUID `json:"team_id,omitempty"`
	Name           string     `json:"name" validate:"required"`
	Description    *string    `json:"description,omitempty"`
	Key            *string    `json:"key,omitempty" validate:"omitempty,max=10"`
	Color          *string    `json:"color,omitempty"`
	Icon           *string    `json:"icon,omitempty"`
	Visibility     *string    `json:"visibility,omitempty" validate:"omitempty,oneof=private team organization"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
}

type UpdateProjectRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Color       *string    `json:"color,omitempty"`
	Icon        *string    `json:"icon,omitempty"`
	Visibility  *string    `json:"visibility,omitempty" validate:"omitempty,oneof=private team organization"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=active on_hold archived completed"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type AddProjectMemberRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Role   string    `json:"role" validate:"required,oneof=owner editor viewer"`
}

// Task DTOs

type TaskDTO struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"project_id"`
	ParentTaskID   *string    `json:"parent_task_id"`
	Title          string     `json:"title"`
	Description    *string    `json:"description"`
	StatusID       string     `json:"status_id"`
	Priority       string     `json:"priority"`
	AssigneeID     *string    `json:"assignee_id"`
	StartDate      *time.Time `json:"start_date"`
	DueDate        *time.Time `json:"due_date"`
	EstimatedHours *float64   `json:"estimated_hours"`
	CustomFields   *string    `json:"custom_fields"`
}
type CreateTaskRequest struct {
	ProjectID      uuid.UUID  `json:"project_id" validate:"required"`
	ParentTaskID   *uuid.UUID `json:"parent_task_id,omitempty"`
	Title          string     `json:"title" validate:"required"`
	Description    *string    `json:"description,omitempty"`
	StatusID       uuid.UUID  `json:"status_id" validate:"required"`
	Priority       *string    `json:"priority,omitempty" validate:"omitempty,oneof=urgent high medium low none"`
	AssigneeID     *uuid.UUID `json:"assignee_id,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedHours *float64   `json:"estimated_hours,omitempty"`
	CustomFields   *string    `json:"custom_fields,omitempty"`
}

type UpdateTaskRequest struct {
	Title          *string    `json:"title,omitempty"`
	Description    *string    `json:"description,omitempty"`
	StatusID       *uuid.UUID `json:"status_id,omitempty"`
	Priority       *string    `json:"priority,omitempty" validate:"omitempty,oneof=urgent high medium low none"`
	AssigneeID     *uuid.UUID `json:"assignee_id,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedHours *float64   `json:"estimated_hours,omitempty"`
	Position       *float64   `json:"position,omitempty"`
}

type TaskFilters struct {
	StatusIDs   []uuid.UUID `json:"status_id,omitempty"`
	AssigneeIDs []uuid.UUID `json:"assignee_id,omitempty"`
	Priorities  []string    `json:"priority,omitempty"`
	Labels      []uuid.UUID `json:"labels,omitempty"`
	Search      string      `json:"search,omitempty"`
	DueDateFrom *time.Time  `json:"due_date_from,omitempty"`
	DueDateTo   *time.Time  `json:"due_date_to,omitempty"`
}

// TaskStatus DTOs
type CreateTaskStatusRequest struct {
	Name     string `json:"name" validate:"required"`
	Color    string `json:"color" validate:"required,hexcolor"`
	Position int    `json:"position" validate:"required"`
}

type UpdateTaskStatusRequest struct {
	Name        *string `json:"name,omitempty"`
	Color       *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Position    *int    `json:"position,omitempty"`
	IsDefault   *bool   `json:"is_default,omitempty"`
	IsCompleted *bool   `json:"is_completed,omitempty"`
}

// Comment DTOs
type CreateCommentRequest struct {
	Content         string     `json:"content" validate:"required"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id,omitempty"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required"`
}

// TimeEntry DTOs
type CreateTimeEntryRequest struct {
	Hours       float64    `json:"hours" validate:"required,gt=0"`
	Description *string    `json:"description,omitempty"`
	StartedAt   time.Time  `json:"started_at" validate:"required"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`
	IsBillable  *bool      `json:"is_billable,omitempty"`
}

// Label DTOs
type CreateLabelRequest struct {
	Name        string     `json:"name" validate:"required"`
	Color       string     `json:"color" validate:"required,hexcolor"`
	Description *string    `json:"description,omitempty"`
	ProjectID   *uuid.UUID `json:"project_id,omitempty"`
}

type UpdateLabelRequest struct {
	Name        *string `json:"name,omitempty"`
	Color       *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Description *string `json:"description,omitempty"`
}

// Team DTOs
type CreateTeamRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
}

type UpdateTeamRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
}

type AddTeamMemberRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// WebSocket message types
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type WSTaskUpdate struct {
	ProjectID uuid.UUID `json:"project_id"`
	Task      *Task     `json:"task"`
	Action    string    `json:"action"` // created, updated, deleted
}

type WSCommentUpdate struct {
	TaskID  uuid.UUID `json:"task_id"`
	Comment *Comment  `json:"comment"`
	Action  string    `json:"action"` // created, updated, deleted
}

// Pagination
type PaginationParams struct {
	Page     int `json:"page" query:"page"`
	PageSize int `json:"page_size" query:"page_size"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

// Generic response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func TransformTaskToDTO(task Task) *TaskDTO {
	dto := &TaskDTO{
		ID:        task.ID.String(),
		Title:     task.Title,
		StatusID:  task.StatusID.String(),
		ProjectID: task.ProjectID.String(),
		Priority:  task.Priority,
	}
	if task.Description.Valid {
		dto.Description = &task.Description.String
	}
	if task.ParentTaskID.Valid {
		val := task.ParentTaskID.UUID.String()
		dto.ParentTaskID = &val
	}
	if task.AssigneeID.Valid {
		val := task.AssigneeID.UUID.String()
		dto.AssigneeID = &val
	}
	if task.StartDate.Valid {
		dto.StartDate = &task.StartDate.Time
	}
	if task.DueDate.Valid {
		dto.DueDate = &task.DueDate.Time
	}
	if task.EstimatedHours.Valid {
		dto.EstimatedHours = &task.EstimatedHours.Float64
	}
	if len(task.CustomFields) > 0 {
		val := string(task.CustomFields)
		dto.CustomFields = &val
	}
	return dto
}

func TransformProjectToDTO(project Project) *ProjectDTO {
	dto := &ProjectDTO{
		ID:         project.ID.String(),
		Name:       project.Name,
		Visibility: project.Visibility,
	}
	if project.Description.Valid {
		dto.Description = &project.Description.String
	}
	if project.Color.Valid {
		dto.Color = &project.Color.String
	}
	if project.Icon.Valid {
		dto.Icon = &project.Icon.String
	}
	if project.StartDate.Valid {
		dto.StartDate = &project.StartDate.Time
	}
	if project.DueDate.Valid {
		dto.DueDate = &project.DueDate.Time
	}
	return dto
}
