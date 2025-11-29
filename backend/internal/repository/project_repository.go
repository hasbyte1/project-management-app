package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hasbyte1/project-management-app/internal/models"
	"github.com/hasbyte1/project-management-app/pkg/database"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	List(ctx context.Context, organizationID uuid.UUID) ([]models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	Archive(ctx context.Context, id uuid.UUID) error
	Unarchive(ctx context.Context, id uuid.UUID) error
	GetMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error)
	AddMember(ctx context.Context, member *models.ProjectMember) error
	UpdateMemberRole(ctx context.Context, id uuid.UUID, role string) error
	RemoveMember(ctx context.Context, id uuid.UUID) error
}

type projectRepository struct {
	db *database.DB
}

func NewProjectRepository(db *database.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO projects (
			id, organization_id, team_id, name, description, key,
			color, icon, visibility, status, start_date, due_date,
			settings, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		project.ID, project.OrganizationID, project.TeamID, project.Name,
		project.Description, project.Key, project.Color, project.Icon,
		project.Visibility, project.Status, project.StartDate, project.DueDate,
		project.Settings, project.CreatedBy,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)

	return err
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	query := `
		SELECT id, organization_id, team_id, name, description, key,
			   color, icon, visibility, status, start_date, due_date,
			   settings, created_by, created_at, updated_at, archived_at, deleted_at
		FROM projects
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &project, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) List(ctx context.Context, organizationID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	query := `
		SELECT id, organization_id, team_id, name, description, key,
			   color, icon, visibility, status, start_date, due_date,
			   settings, created_by, created_at, updated_at, archived_at, deleted_at
		FROM projects
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &projects, query, organizationID)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, project *models.Project) error {
	query := `
		UPDATE projects
		SET name = $1, description = $2, color = $3, icon = $4,
			visibility = $5, status = $6, start_date = $7, due_date = $8,
			updated_at = NOW()
		WHERE id = $9 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx, query,
		project.Name, project.Description, project.Color, project.Icon,
		project.Visibility, project.Status, project.StartDate, project.DueDate,
		project.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE projects SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *projectRepository) Archive(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE projects
		SET archived_at = NOW(), status = 'archived', updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *projectRepository) Unarchive(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE projects
		SET archived_at = NULL, status = 'active', updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *projectRepository) GetMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error) {
	var members []models.ProjectMember
	query := `
		SELECT pm.id, pm.project_id, pm.user_id, pm.role, pm.added_by,
			   pm.created_at, pm.updated_at,
			   u.id as "user.id", u.email as "user.email",
			   u.first_name as "user.first_name", u.last_name as "user.last_name",
			   u.avatar_url as "user.avatar_url"
		FROM project_members pm
		JOIN users u ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY pm.created_at ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member models.ProjectMember
		member.User = &models.User{}

		err := rows.Scan(
			&member.ID, &member.ProjectID, &member.UserID, &member.Role, &member.AddedBy,
			&member.CreatedAt, &member.UpdatedAt,
			&member.User.ID, &member.User.Email, &member.User.FirstName,
			&member.User.LastName, &member.User.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *projectRepository) AddMember(ctx context.Context, member *models.ProjectMember) error {
	query := `
		INSERT INTO project_members (id, project_id, user_id, role, added_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		member.ID, member.ProjectID, member.UserID, member.Role, member.AddedBy,
	).Scan(&member.ID, &member.CreatedAt, &member.UpdatedAt)

	return err
}

func (r *projectRepository) UpdateMemberRole(ctx context.Context, id uuid.UUID, role string) error {
	query := `
		UPDATE project_members
		SET role = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, role, id)
	return err
}

func (r *projectRepository) RemoveMember(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM project_members WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
