package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hasbyte1/project-management-app/internal/models"
	"github.com/hasbyte1/project-management-app/pkg/database"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error)
	List(ctx context.Context, projectID uuid.UUID, filters *models.TaskFilters) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID) error
	UpdatePosition(ctx context.Context, id uuid.UUID, position float64) error
	GetNextTaskNumber(ctx context.Context, projectID uuid.UUID) (int, error)

	// Status management
	GetStatuses(ctx context.Context, projectID uuid.UUID) ([]models.TaskStatus, error)
	CreateStatus(ctx context.Context, status *models.TaskStatus) error
	UpdateTaskStatus(ctx context.Context, status *models.TaskStatus) error
	DeleteStatus(ctx context.Context, id uuid.UUID) error

	// Comments
	GetComments(ctx context.Context, taskID uuid.UUID) ([]models.Comment, error)
	CreateComment(ctx context.Context, comment *models.Comment) error
	UpdateComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, id uuid.UUID) error

	// Labels
	GetLabels(ctx context.Context, taskID uuid.UUID) ([]models.Label, error)
	AddLabel(ctx context.Context, taskID, labelID uuid.UUID) error
	RemoveLabel(ctx context.Context, taskID, labelID uuid.UUID) error
}

type taskRepository struct {
	db *database.DB
}

func NewTaskRepository(db *database.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO tasks (
			id, project_id, parent_task_id, title, description, task_number,
			status_id, priority, assignee_id, reporter_id, start_date, due_date,
			estimated_hours, position, custom_fields, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		task.ID, task.ProjectID, task.ParentTaskID, task.Title, task.Description,
		task.TaskNumber, task.StatusID, task.Priority, task.AssigneeID, task.ReporterID,
		task.StartDate, task.DueDate, task.EstimatedHours, task.Position,
		task.CustomFields, task.CreatedBy,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	return err
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	var task models.Task
	query := `
		SELECT t.id, t.project_id, t.parent_task_id, t.title, t.description,
			   t.task_number, t.status_id, t.priority, t.assignee_id, t.reporter_id,
			   t.start_date, t.due_date, t.completed_at, t.estimated_hours,
			   t.actual_hours, t.position, t.custom_fields, t.created_by,
			   t.created_at, t.updated_at, t.deleted_at,
			   ts.id as "status.id", ts.name as "status.name", ts.color as "status.color",
			   ts.position as "status.position", ts.is_default as "status.is_default",
			   ts.is_completed as "status.is_completed"
		FROM tasks t
		JOIN task_statuses ts ON ts.id = t.status_id
		WHERE t.id = $1 AND t.deleted_at IS NULL
	`

	row := r.db.QueryRowxContext(ctx, query, id)
	task.Status = &models.TaskStatus{}

	err := row.Scan(
		&task.ID, &task.ProjectID, &task.ParentTaskID, &task.Title, &task.Description,
		&task.TaskNumber, &task.StatusID, &task.Priority, &task.AssigneeID, &task.ReporterID,
		&task.StartDate, &task.DueDate, &task.CompletedAt, &task.EstimatedHours,
		&task.ActualHours, &task.Position, &task.CustomFields, &task.CreatedBy,
		&task.CreatedAt, &task.UpdatedAt, &task.DeletedAt,
		&task.Status.ID, &task.Status.Name, &task.Status.Color,
		&task.Status.Position, &task.Status.IsDefault, &task.Status.IsCompleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	// Load relations
	if task.AssigneeID.Valid {
		// Load assignee separately
	}

	// Load labels
	labels, _ := r.GetLabels(ctx, task.ID)
	task.Labels = labels

	return &task, nil
}

func (r *taskRepository) List(ctx context.Context, projectID uuid.UUID, filters *models.TaskFilters) ([]models.Task, error) {
	query := `
		SELECT t.id, t.project_id, t.parent_task_id, t.title, t.description,
			   t.task_number, t.status_id, t.priority, t.assignee_id, t.reporter_id,
			   t.start_date, t.due_date, t.completed_at, t.estimated_hours,
			   t.actual_hours, t.position, t.custom_fields, t.created_by,
			   t.created_at, t.updated_at,
			   ts.id as "status.id", ts.name as "status.name", ts.color as "status.color"
		FROM tasks t
		JOIN task_statuses ts ON ts.id = t.status_id
		WHERE t.project_id = $1 AND t.deleted_at IS NULL
	`

	args := []interface{}{projectID}
	argPos := 2

	// Apply filters
	if filters != nil {
		if len(filters.StatusIDs) > 0 {
			placeholders := make([]string, len(filters.StatusIDs))
			for i, statusID := range filters.StatusIDs {
				placeholders[i] = fmt.Sprintf("$%d", argPos)
				args = append(args, statusID)
				argPos++
			}
			query += fmt.Sprintf(" AND t.status_id IN (%s)", strings.Join(placeholders, ","))
		}

		if len(filters.AssigneeIDs) > 0 {
			placeholders := make([]string, len(filters.AssigneeIDs))
			for i, assigneeID := range filters.AssigneeIDs {
				placeholders[i] = fmt.Sprintf("$%d", argPos)
				args = append(args, assigneeID)
				argPos++
			}
			query += fmt.Sprintf(" AND t.assignee_id IN (%s)", strings.Join(placeholders, ","))
		}

		if len(filters.Priorities) > 0 {
			placeholders := make([]string, len(filters.Priorities))
			for i, priority := range filters.Priorities {
				placeholders[i] = fmt.Sprintf("$%d", argPos)
				args = append(args, priority)
				argPos++
			}
			query += fmt.Sprintf(" AND t.priority IN (%s)", strings.Join(placeholders, ","))
		}

		if filters.Search != "" {
			query += fmt.Sprintf(" AND (t.title ILIKE $%d OR t.description ILIKE $%d)", argPos, argPos)
			searchTerm := "%" + filters.Search + "%"
			args = append(args, searchTerm)
			argPos++
		}
	}

	query += " ORDER BY t.position ASC, t.created_at DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		task.Status = &models.TaskStatus{}

		err := rows.Scan(
			&task.ID, &task.ProjectID, &task.ParentTaskID, &task.Title, &task.Description,
			&task.TaskNumber, &task.StatusID, &task.Priority, &task.AssigneeID, &task.ReporterID,
			&task.StartDate, &task.DueDate, &task.CompletedAt, &task.EstimatedHours,
			&task.ActualHours, &task.Position, &task.CustomFields, &task.CreatedBy,
			&task.CreatedAt, &task.UpdatedAt,
			&task.Status.ID, &task.Status.Name, &task.Status.Color,
		)
		if err != nil {
			return nil, err
		}

		// Load labels for each task
		labels, _ := r.GetLabels(ctx, task.ID)
		task.Labels = labels

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status_id = $3, priority = $4,
			assignee_id = $5, start_date = $6, due_date = $7,
			estimated_hours = $8, position = $9, updated_at = NOW()
		WHERE id = $10 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(
		ctx, query,
		task.Title, task.Description, task.StatusID, task.Priority,
		task.AssigneeID, task.StartDate, task.DueDate,
		task.EstimatedHours, task.Position, task.ID,
	)

	return err
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tasks SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *taskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, statusID uuid.UUID) error {
	query := `
		UPDATE tasks
		SET status_id = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, statusID, id)
	return err
}

func (r *taskRepository) UpdatePosition(ctx context.Context, id uuid.UUID, position float64) error {
	query := `
		UPDATE tasks
		SET position = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, position, id)
	return err
}

func (r *taskRepository) GetNextTaskNumber(ctx context.Context, projectID uuid.UUID) (int, error) {
	var maxNumber sql.NullInt32
	query := `
		SELECT MAX(task_number)
		FROM tasks
		WHERE project_id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRowContext(ctx, query, projectID).Scan(&maxNumber)
	if err != nil {
		return 0, err
	}

	if maxNumber.Valid {
		return int(maxNumber.Int32) + 1, nil
	}

	return 1, nil
}

// Task Statuses
func (r *taskRepository) GetStatuses(ctx context.Context, projectID uuid.UUID) ([]models.TaskStatus, error) {
	var statuses []models.TaskStatus
	query := `
		SELECT id, project_id, name, color, position, is_default, is_completed,
			   created_at, updated_at
		FROM task_statuses
		WHERE project_id = $1
		ORDER BY position ASC
	`

	err := r.db.SelectContext(ctx, &statuses, query, projectID)
	return statuses, err
}

func (r *taskRepository) CreateStatus(ctx context.Context, status *models.TaskStatus) error {
	query := `
		INSERT INTO task_statuses (id, project_id, name, color, position, is_default, is_completed)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		status.ID, status.ProjectID, status.Name, status.Color,
		status.Position, status.IsDefault, status.IsCompleted,
	).Scan(&status.ID, &status.CreatedAt, &status.UpdatedAt)

	return err
}

func (r *taskRepository) UpdateTaskStatus(ctx context.Context, status *models.TaskStatus) error {
	query := `
		UPDATE task_statuses
		SET name = $1, color = $2, position = $3, is_default = $4,
			is_completed = $5, updated_at = NOW()
		WHERE id = $6
	`

	_, err := r.db.ExecContext(
		ctx, query,
		status.Name, status.Color, status.Position,
		status.IsDefault, status.IsCompleted, status.ID,
	)

	return err
}

func (r *taskRepository) DeleteStatus(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM task_statuses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Comments
func (r *taskRepository) GetComments(ctx context.Context, taskID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	query := `
		SELECT c.id, c.task_id, c.user_id, c.content, c.parent_comment_id,
			   c.is_edited, c.created_at, c.updated_at,
			   u.id as "user.id", u.email as "user.email",
			   u.first_name as "user.first_name", u.last_name as "user.last_name",
			   u.avatar_url as "user.avatar_url"
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.task_id = $1 AND c.deleted_at IS NULL
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		comment.User = &models.User{}

		err := rows.Scan(
			&comment.ID, &comment.TaskID, &comment.UserID, &comment.Content,
			&comment.ParentCommentID, &comment.IsEdited, &comment.CreatedAt, &comment.UpdatedAt,
			&comment.User.ID, &comment.User.Email, &comment.User.FirstName,
			&comment.User.LastName, &comment.User.AvatarURL,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *taskRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (id, task_id, user_id, content, parent_comment_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		comment.ID, comment.TaskID, comment.UserID, comment.Content, comment.ParentCommentID,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)

	return err
}

func (r *taskRepository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	query := `
		UPDATE comments
		SET content = $1, is_edited = true, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, comment.Content, comment.ID)
	return err
}

func (r *taskRepository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE comments SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Labels
func (r *taskRepository) GetLabels(ctx context.Context, taskID uuid.UUID) ([]models.Label, error) {
	var labels []models.Label
	query := `
		SELECT l.id, l.organization_id, l.project_id, l.name, l.color,
			   l.description, l.created_by, l.created_at, l.updated_at
		FROM labels l
		JOIN task_labels tl ON tl.label_id = l.id
		WHERE tl.task_id = $1 AND l.deleted_at IS NULL
	`

	err := r.db.SelectContext(ctx, &labels, query, taskID)
	return labels, err
}

func (r *taskRepository) AddLabel(ctx context.Context, taskID, labelID uuid.UUID) error {
	query := `
		INSERT INTO task_labels (id, task_id, label_id)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, uuid.New(), taskID, labelID)
	return err
}

func (r *taskRepository) RemoveLabel(ctx context.Context, taskID, labelID uuid.UUID) error {
	query := `DELETE FROM task_labels WHERE task_id = $1 AND label_id = $2`
	_, err := r.db.ExecContext(ctx, query, taskID, labelID)
	return err
}
