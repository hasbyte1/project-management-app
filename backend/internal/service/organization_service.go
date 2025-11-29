package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hasbyte1/project-management-app/internal/models"
	"github.com/hasbyte1/project-management-app/internal/repository"
)

type OrganizationService interface {
	Create(ctx context.Context, req *models.CreateOrganizationRequest, userID uuid.UUID) (*models.Organization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	List(ctx context.Context, userID uuid.UUID) ([]models.OrganizationDTO, error)
	Update(ctx context.Context, id uuid.UUID, req *models.UpdateOrganizationRequest) (*models.Organization, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMember, error)
	AddMember(ctx context.Context, orgID uuid.UUID, req *models.AddOrganizationMemberRequest, invitedBy uuid.UUID) (*models.OrganizationMember, error)
	UpdateMemberRole(ctx context.Context, memberID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, memberID uuid.UUID) error
}

type organizationService struct {
	orgRepo  repository.OrganizationRepository
	userRepo repository.UserRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository, userRepo repository.UserRepository) OrganizationService {
	return &organizationService{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

func (s *organizationService) Create(ctx context.Context, req *models.CreateOrganizationRequest, userID uuid.UUID) (*models.Organization, error) {
	// Check if slug is unique
	existing, _ := s.orgRepo.GetBySlug(ctx, req.Slug)
	if existing != nil {
		return nil, fmt.Errorf("organization with slug %s already exists", req.Slug)
	}

	// Calculate depth and path
	depth := 0
	var path sql.NullString

	if req.ParentID != nil {
		parent, err := s.orgRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent organization not found")
		}
		depth = parent.Depth + 1
		if parent.Path.Valid {
			path.String = fmt.Sprintf("%s.%s", parent.Path.String, uuid.New().String())
		} else {
			path.String = uuid.New().String()
		}
		path.Valid = true
	} else {
		path.String = uuid.New().String()
		path.Valid = true
	}

	org := &models.Organization{
		Base: models.Base{
			ID: uuid.New(),
		},
		Name:      req.Name,
		Slug:      req.Slug,
		Depth:     depth,
		Path:      path,
		Settings:  []byte("{}"),
		CreatedBy: userID,
	}

	if req.ParentID != nil {
		org.ParentID = uuid.NullUUID{UUID: *req.ParentID, Valid: true}
	}

	if req.Description != nil {
		org.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.LogoURL != nil {
		org.LogoURL = sql.NullString{String: *req.LogoURL, Valid: true}
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Add creator as owner
	member := &models.OrganizationMember{
		Base: models.Base{
			ID: uuid.New(),
		},
		OrganizationID: org.ID,
		UserID:         userID,
		Role:           "owner",
		JoinedAt:       time.Now(),
	}

	if err := s.orgRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add creator as member: %w", err)
	}

	return org, nil
}

func (s *organizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	return s.orgRepo.GetByID(ctx, id)
}

func (s *organizationService) List(ctx context.Context, userID uuid.UUID) ([]models.OrganizationDTO, error) {
	orgs, err := s.orgRepo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	var orgDTOs []models.OrganizationDTO
	for _, org := range orgs {
		row := models.OrganizationDTO{
			ID:        org.ID.String(),
			Name:      org.Name,
			Slug:      org.Slug,
			Depth:     org.Depth,
			Settings:  org.Settings,
			CreatedBy: org.CreatedBy.String(),
		}
		orgDTOs = append(orgDTOs, row)

		if org.ParentID.Valid {
			val := org.ParentID.UUID.String()
			row.ParentID = &val
		}
		if org.Description.Valid {
			row.Description = &org.Description.String
		}
		if org.Path.Valid {
			row.Path = &org.Path.String
		}
		if org.Path.Valid {
			row.Path = &org.Path.String
		}
	}

	return orgDTOs, nil
}

func (s *organizationService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateOrganizationRequest) (*models.Organization, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		org.Name = *req.Name
	}

	if req.Description != nil {
		org.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.LogoURL != nil {
		org.LogoURL = sql.NullString{String: *req.LogoURL, Valid: true}
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}

	return s.orgRepo.GetByID(ctx, id)
}

func (s *organizationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.orgRepo.Delete(ctx, id)
}

func (s *organizationService) GetMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMember, error) {
	return s.orgRepo.GetMembers(ctx, orgID)
}

func (s *organizationService) AddMember(ctx context.Context, orgID uuid.UUID, req *models.AddOrganizationMemberRequest, invitedBy uuid.UUID) (*models.OrganizationMember, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("user with email %s not found", req.Email)
	}

	// Check if already a member
	members, err := s.orgRepo.GetMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		if member.UserID == user.ID {
			return nil, fmt.Errorf("user is already a member of this organization")
		}
	}

	now := time.Now()
	member := &models.OrganizationMember{
		Base: models.Base{
			ID: uuid.New(),
		},
		OrganizationID: orgID,
		UserID:         user.ID,
		Role:           req.Role,
		InvitedBy:      uuid.NullUUID{UUID: invitedBy, Valid: true},
		InvitedAt:      sql.NullTime{Time: now, Valid: true},
		JoinedAt:       now,
	}

	if err := s.orgRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	// Load user details
	member.User = user

	return member, nil
}

func (s *organizationService) UpdateMemberRole(ctx context.Context, memberID uuid.UUID, role string) error {
	return s.orgRepo.UpdateMemberRole(ctx, memberID, role)
}

func (s *organizationService) RemoveMember(ctx context.Context, memberID uuid.UUID) error {
	return s.orgRepo.RemoveMember(ctx, memberID)
}
