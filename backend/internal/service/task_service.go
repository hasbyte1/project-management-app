package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hasbyte1/project-management-app/internal/models"
	"github.com/hasbyte1/project-management-app/internal/repository"
)

type TaskService interface {
	Create(ctx context.Context, req *models.CreateTaskRequest, userID uuid.UUID) (*models.TaskDTO, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.TaskDTO, error)
	List(ctx context.Context, projectID uuid.UUID, filters *models.TaskFilters) ([]models.TaskDTO, error)
	Update(ctx context.Context, id uuid.UUID, req *models.UpdateTaskRequest) (*models.TaskDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID) (*models.TaskDTO, error)
	GetStatuses(ctx context.Context, projectID uuid.UUID) ([]models.TaskStatusDTO, error)
	CreateStatus(ctx context.Context, projectID uuid.UUID, req *models.CreateTaskStatusRequest) (*models.TaskStatusDTO, error)
	GetComments(ctx context.Context, taskID uuid.UUID) ([]models.Comment, error)
	CreateComment(ctx context.Context, taskID uuid.UUID, req *models.CreateCommentRequest, userID uuid.UUID) (*models.Comment, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) Create(ctx context.Context, req *models.CreateTaskRequest, userID uuid.UUID) (*models.TaskDTO, error) {
	// Get next task number
	taskNumber, err := s.taskRepo.GetNextTaskNumber(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}

	task := &models.Task{
		Base:         models.Base{ID: uuid.New()},
		ProjectID:    req.ProjectID,
		ParentTaskID: uuid.NullUUID{Valid: req.ParentTaskID != nil},
		Title:        req.Title,
		TaskNumber:   taskNumber,
		StatusID:     req.StatusID,
		Priority:     "none",
		ReporterID:   userID,
		CreatedBy:    userID,
		ActualHours:  0,
		Position:     float64(taskNumber),
		CustomFields: []byte("{}"),
	}

	if req.ParentTaskID != nil {
		task.ParentTaskID.UUID = *req.ParentTaskID
	}

	if req.Description != nil {
		task.Description.String = *req.Description
		task.Description.Valid = true
	}

	if req.Priority != nil {
		task.Priority = *req.Priority
	}

	if req.AssigneeID != nil {
		task.AssigneeID = uuid.NullUUID{UUID: *req.AssigneeID, Valid: true}
	}

	if req.StartDate != nil {
		task.StartDate.Time = *req.StartDate
		task.StartDate.Valid = true
	}

	if req.DueDate != nil {
		task.DueDate.Time = *req.DueDate
		task.DueDate.Valid = true
	}

	if req.EstimatedHours != nil {
		task.EstimatedHours.Float64 = *req.EstimatedHours
		task.EstimatedHours.Valid = true
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	record, err := s.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	dto := models.TransformTaskToDTO(*record)
	return &dto, nil
}

func (s *taskService) GetByID(ctx context.Context, id uuid.UUID) (*models.TaskDTO, error) {
	record, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	dto := models.TransformTaskToDTO(*record)

	return &dto, nil
}

func (s *taskService) List(ctx context.Context, projectID uuid.UUID, filters *models.TaskFilters) ([]models.TaskDTO, error) {
	rows, err := s.taskRepo.List(ctx, projectID, filters)
	if err != nil {
		return nil, err
	}

	var tasks []models.TaskDTO
	for _, row := range rows {
		tasks = append(tasks, models.TransformTaskToDTO(row))
	}

	return tasks, nil
}

func (s *taskService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateTaskRequest) (*models.TaskDTO, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		task.Title = *req.Title
	}

	if req.Description != nil {
		task.Description.String = *req.Description
		task.Description.Valid = true
	}

	if req.StatusID != nil {
		task.StatusID = *req.StatusID
	}

	if req.Priority != nil {
		task.Priority = *req.Priority
	}

	if req.AssigneeID != nil {
		task.AssigneeID = uuid.NullUUID{UUID: *req.AssigneeID, Valid: true}
	}

	if req.StartDate != nil {
		task.StartDate.Time = *req.StartDate
		task.StartDate.Valid = true
	}

	if req.DueDate != nil {
		task.DueDate.Time = *req.DueDate
		task.DueDate.Valid = true
	}

	if req.EstimatedHours != nil {
		task.EstimatedHours.Float64 = *req.EstimatedHours
		task.EstimatedHours.Valid = true
	}

	if req.Position != nil {
		task.Position = *req.Position
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	record, err := s.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	dto := models.TransformTaskToDTO(*record)

	return &dto, nil
}

func (s *taskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *taskService) UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID) (*models.TaskDTO, error) {
	if err := s.taskRepo.UpdateStatus(ctx, id, statusID); err != nil {
		return nil, err
	}
	record, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	dto := models.TransformTaskToDTO(*record)
	return &dto, nil
}

func (s *taskService) GetStatuses(ctx context.Context, projectID uuid.UUID) ([]models.TaskStatusDTO, error) {
	records, err := s.taskRepo.GetStatuses(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var statuses []models.TaskStatusDTO
	for _, record := range records {
		statuses = append(statuses, models.TransformTaskStatusToDTO(record))
	}

	return statuses, nil
}

func (s *taskService) CreateStatus(ctx context.Context, projectID uuid.UUID, req *models.CreateTaskStatusRequest) (*models.TaskStatusDTO, error) {
	status := &models.TaskStatus{
		Base:        models.Base{ID: uuid.New()},
		ProjectID:   projectID,
		Name:        req.Name,
		Color:       req.Color,
		Position:    req.Position,
		IsDefault:   false,
		IsCompleted: false,
	}

	if err := s.taskRepo.CreateStatus(ctx, status); err != nil {
		return nil, err
	}

	dto := models.TransformTaskStatusToDTO(*status)

	return &dto, nil
}

func (s *taskService) GetComments(ctx context.Context, taskID uuid.UUID) ([]models.Comment, error) {
	return s.taskRepo.GetComments(ctx, taskID)
}

func (s *taskService) CreateComment(ctx context.Context, taskID uuid.UUID, req *models.CreateCommentRequest, userID uuid.UUID) (*models.Comment, error) {
	comment := &models.Comment{
		Base:    models.Base{ID: uuid.New()},
		TaskID:  taskID,
		UserID:  userID,
		Content: req.Content,
	}

	if req.ParentCommentID != nil {
		comment.ParentCommentID = uuid.NullUUID{UUID: *req.ParentCommentID, Valid: true}
	}

	if err := s.taskRepo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}
