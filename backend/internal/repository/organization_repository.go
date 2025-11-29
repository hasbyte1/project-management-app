package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hasbyte1/project-management-app/internal/models"
	"github.com/hasbyte1/project-management-app/pkg/database"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *models.Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	GetBySlug(ctx context.Context, slug string) (*models.Organization, error)
	List(ctx context.Context, userID uuid.UUID) ([]models.Organization, error)
	Update(ctx context.Context, org *models.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMember, error)
	AddMember(ctx context.Context, member *models.OrganizationMember) error
	UpdateMemberRole(ctx context.Context, id uuid.UUID, role string) error
	RemoveMember(ctx context.Context, id uuid.UUID) error
}

type organizationRepository struct {
	db *database.DB
}

func NewOrganizationRepository(db *database.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, org *models.Organization) error {
	query := `
		INSERT INTO organizations (
			id, parent_id, name, slug, description, logo_url,
			depth, path, settings, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		org.ID, org.ParentID, org.Name, org.Slug, org.Description,
		org.LogoURL, org.Depth, org.Path, org.Settings, org.CreatedBy,
	).Scan(&org.ID, &org.CreatedAt, &org.UpdatedAt)

	return err
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	query := `
		SELECT id, parent_id, name, slug, description, logo_url,
			   depth, path, settings, created_by,
			   created_at, updated_at, deleted_at
		FROM organizations
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &org, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, err
	}

	return &org, nil
}

func (r *organizationRepository) GetBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	var org models.Organization
	query := `
		SELECT id, parent_id, name, slug, description, logo_url,
			   depth, path, settings, created_by,
			   created_at, updated_at, deleted_at
		FROM organizations
		WHERE slug = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &org, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, err
	}

	return &org, nil
}

func (r *organizationRepository) List(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	var orgs []models.Organization
	query := `
		SELECT DISTINCT o.id, o.parent_id, o.name, o.slug, o.description,
			   o.logo_url, o.depth, o.path, o.settings, o.created_by,
			   o.created_at, o.updated_at
		FROM organizations o
		JOIN organization_members om ON om.organization_id = o.id
		WHERE om.user_id = $1 AND o.deleted_at IS NULL
		ORDER BY o.created_at DESC
	`

	err := r.db.SelectContext(ctx, &orgs, query, userID)
	return orgs, err
}

func (r *organizationRepository) Update(ctx context.Context, org *models.Organization) error {
	query := `
		UPDATE organizations
		SET name = $1, description = $2, logo_url = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(
		ctx, query,
		org.Name, org.Description, org.LogoURL, org.ID,
	)

	return err
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE organizations SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *organizationRepository) GetMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMember, error) {
	var members []models.OrganizationMember
	query := `
		SELECT om.id, om.organization_id, om.user_id, om.role, om.invited_by,
			   om.invited_at, om.joined_at, om.created_at, om.updated_at,
			   u.id as "user.id", u.email as "user.email",
			   u.first_name as "user.first_name", u.last_name as "user.last_name",
			   u.avatar_url as "user.avatar_url"
		FROM organization_members om
		JOIN users u ON u.id = om.user_id
		WHERE om.organization_id = $1
		ORDER BY om.created_at ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member models.OrganizationMember
		member.User = &models.User{}

		err := rows.Scan(
			&member.ID, &member.OrganizationID, &member.UserID, &member.Role,
			&member.InvitedBy, &member.InvitedAt, &member.JoinedAt,
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

func (r *organizationRepository) AddMember(ctx context.Context, member *models.OrganizationMember) error {
	query := `
		INSERT INTO organization_members (
			id, organization_id, user_id, role, invited_by, invited_at, joined_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		member.ID, member.OrganizationID, member.UserID, member.Role,
		member.InvitedBy, member.InvitedAt, member.JoinedAt,
	).Scan(&member.ID, &member.CreatedAt, &member.UpdatedAt)

	return err
}

func (r *organizationRepository) UpdateMemberRole(ctx context.Context, id uuid.UUID, role string) error {
	query := `
		UPDATE organization_members
		SET role = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, role, id)
	return err
}

func (r *organizationRepository) RemoveMember(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM organization_members WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
