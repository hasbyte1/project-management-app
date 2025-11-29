package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hasbyte1/project-management-app/internal/models"
	"github.com/hasbyte1/project-management-app/internal/repository"
)

type ProjectService interface {
	Create(ctx context.Context, req *models.CreateProjectRequest, userID uuid.UUID) (*models.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	List(ctx context.Context, organizationID uuid.UUID) ([]models.ProjectDTO, error)
	Update(ctx context.Context, id uuid.UUID, req *models.UpdateProjectRequest) (*models.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Archive(ctx context.Context, id uuid.UUID) error
	Unarchive(ctx context.Context, id uuid.UUID) error
	GetMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error)
	AddMember(ctx context.Context, projectID uuid.UUID, req *models.AddProjectMemberRequest, addedBy uuid.UUID) (*models.ProjectMember, error)
	UpdateMemberRole(ctx context.Context, memberID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, memberID uuid.UUID) error
}

type projectService struct {
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository
}

func NewProjectService(projectRepo repository.ProjectRepository, taskRepo repository.TaskRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
	}
}

func (s *projectService) Create(ctx context.Context, req *models.CreateProjectRequest, userID uuid.UUID) (*models.Project, error) {
	project := &models.Project{
		Base: models.Base{
			ID: uuid.New(),
		},
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Visibility:     "team",
		Status:         "active",
		Settings:       []byte("{}"),
		CreatedBy:      userID,
	}

	if req.TeamID != nil {
		project.TeamID = uuid.NullUUID{UUID: *req.TeamID, Valid: true}
	}

	if req.Description != nil {
		project.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.Key != nil {
		project.Key = sql.NullString{String: *req.Key, Valid: true}
	}

	if req.Color != nil {
		project.Color = sql.NullString{String: *req.Color, Valid: true}
	}

	if req.Icon != nil {
		project.Icon = sql.NullString{String: *req.Icon, Valid: true}
	}

	if req.Visibility != nil {
		project.Visibility = *req.Visibility
	}

	if req.StartDate != nil {
		project.StartDate = sql.NullTime{Time: *req.StartDate, Valid: true}
	}

	if req.DueDate != nil {
		project.DueDate = sql.NullTime{Time: *req.DueDate, Valid: true}
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Add creator as owner
	member := &models.ProjectMember{
		Base: models.Base{
			ID: uuid.New(),
		},
		ProjectID: project.ID,
		UserID:    userID,
		Role:      "owner",
		AddedBy:   uuid.NullUUID{UUID: userID, Valid: true},
	}

	if err := s.projectRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add creator as member: %w", err)
	}

	// Create default task statuses
	if err := s.createDefaultStatuses(ctx, project.ID); err != nil {
		return nil, fmt.Errorf("failed to create default statuses: %w", err)
	}

	return project, nil
}

func (s *projectService) createDefaultStatuses(ctx context.Context, projectID uuid.UUID) error {
	statuses := []struct {
		name        string
		color       string
		position    int
		isDefault   bool
		isCompleted bool
	}{
		{"Backlog", "#94a3b8", 1, false, false},
		{"To Do", "#3b82f6", 2, true, false},
		{"In Progress", "#f59e0b", 3, false, false},
		{"In Review", "#8b5cf6", 4, false, false},
		{"Done", "#10b981", 5, false, true},
	}

	for _, st := range statuses {
		status := &models.TaskStatus{
			Base: models.Base{
				ID: uuid.New(),
			},
			ProjectID:   projectID,
			Name:        st.name,
			Color:       st.color,
			Position:    st.position,
			IsDefault:   st.isDefault,
			IsCompleted: st.isCompleted,
		}

		if err := s.taskRepo.CreateStatus(ctx, status); err != nil {
			return err
		}
	}

	return nil
}

func (s *projectService) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

func (s *projectService) List(ctx context.Context, organizationID uuid.UUID) ([]models.ProjectDTO, error) {
	rows, err := s.projectRepo.List(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	result := []models.ProjectDTO{}
	for _, row := range rows {
		record := models.TransformProjectToDTO(row)
		result = append(result, *record)
	}
	return result, nil
}

func (s *projectService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateProjectRequest) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		project.Name = *req.Name
	}

	if req.Description != nil {
		project.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.Color != nil {
		project.Color = sql.NullString{String: *req.Color, Valid: true}
	}

	if req.Icon != nil {
		project.Icon = sql.NullString{String: *req.Icon, Valid: true}
	}

	if req.Visibility != nil {
		project.Visibility = *req.Visibility
	}

	if req.Status != nil {
		project.Status = *req.Status
	}

	if req.StartDate != nil {
		project.StartDate = sql.NullTime{Time: *req.StartDate, Valid: true}
	}

	if req.DueDate != nil {
		project.DueDate = sql.NullTime{Time: *req.DueDate, Valid: true}
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return s.projectRepo.GetByID(ctx, id)
}

func (s *projectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.projectRepo.Delete(ctx, id)
}

func (s *projectService) Archive(ctx context.Context, id uuid.UUID) error {
	return s.projectRepo.Archive(ctx, id)
}

func (s *projectService) Unarchive(ctx context.Context, id uuid.UUID) error {
	return s.projectRepo.Unarchive(ctx, id)
}

func (s *projectService) GetMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error) {
	return s.projectRepo.GetMembers(ctx, projectID)
}

func (s *projectService) AddMember(ctx context.Context, projectID uuid.UUID, req *models.AddProjectMemberRequest, addedBy uuid.UUID) (*models.ProjectMember, error) {
	// Check if already a member
	members, err := s.projectRepo.GetMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		if member.UserID == req.UserID {
			return nil, fmt.Errorf("user is already a member of this project")
		}
	}

	member := &models.ProjectMember{
		Base: models.Base{
			ID: uuid.New(),
		},
		ProjectID: projectID,
		UserID:    req.UserID,
		Role:      req.Role,
		AddedBy:   uuid.NullUUID{UUID: addedBy, Valid: true},
	}

	if err := s.projectRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *projectService) UpdateMemberRole(ctx context.Context, memberID uuid.UUID, role string) error {
	return s.projectRepo.UpdateMemberRole(ctx, memberID, role)
}

func (s *projectService) RemoveMember(ctx context.Context, memberID uuid.UUID) error {
	return s.projectRepo.RemoveMember(ctx, memberID)
}
