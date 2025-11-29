package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Base model with common fields
type Base struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}

// User represents a system user
type User struct {
	Base
	Email         string         `json:"email" db:"email"`
	PasswordHash  string         `json:"-" db:"password_hash"`
	FirstName     string         `json:"first_name" db:"first_name"`
	LastName      string         `json:"last_name" db:"last_name"`
	AvatarURL     sql.NullString `json:"avatar_url,omitempty" db:"avatar_url"`
	Timezone      string         `json:"timezone" db:"timezone"`
	Locale        string         `json:"locale" db:"locale"`
	EmailVerified bool           `json:"email_verified" db:"email_verified"`
	IsActive      bool           `json:"is_active" db:"is_active"`
	LastLoginAt   sql.NullTime   `json:"last_login_at,omitempty" db:"last_login_at"`
}

// UserToken represents authentication tokens
type UserToken struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	UserID    uuid.UUID    `json:"user_id" db:"user_id"`
	TokenHash string       `json:"-" db:"token_hash"`
	TokenType string       `json:"token_type" db:"token_type"`
	ExpiresAt time.Time    `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UsedAt    sql.NullTime `json:"used_at,omitempty" db:"used_at"`
}

// Organization represents a top-level organization
type Organization struct {
	Base
	ParentID    uuid.NullUUID  `json:"parent_id,omitempty" db:"parent_id"`
	Name        string         `json:"name" db:"name"`
	Slug        string         `json:"slug" db:"slug"`
	Description sql.NullString `json:"description,omitempty" db:"description"`
	LogoURL     sql.NullString `json:"logo_url,omitempty" db:"logo_url"`
	Depth       int            `json:"depth" db:"depth"`
	Path        sql.NullString `json:"path,omitempty" db:"path"`
	Settings    []byte         `json:"settings" db:"settings"`
	CreatedBy   uuid.UUID      `json:"created_by" db:"created_by"`
}

// OrganizationMember represents organization membership
type OrganizationMember struct {
	Base
	OrganizationID uuid.UUID      `json:"organization_id" db:"organization_id"`
	UserID         uuid.UUID      `json:"user_id" db:"user_id"`
	Role           string         `json:"role" db:"role"`
	InvitedBy      uuid.NullUUID  `json:"invited_by,omitempty" db:"invited_by"`
	InvitedAt      sql.NullTime   `json:"invited_at,omitempty" db:"invited_at"`
	JoinedAt       time.Time      `json:"joined_at" db:"joined_at"`
	User           *User          `json:"user,omitempty" db:"-"`
}

// Team represents a team within an organization
type Team struct {
	Base
	OrganizationID uuid.UUID      `json:"organization_id" db:"organization_id"`
	Name           string         `json:"name" db:"name"`
	Description    sql.NullString `json:"description,omitempty" db:"description"`
	Color          sql.NullString `json:"color,omitempty" db:"color"`
	CreatedBy      uuid.UUID      `json:"created_by" db:"created_by"`
}

// TeamMember represents team membership
type TeamMember struct {
	ID        uuid.UUID     `json:"id" db:"id"`
	TeamID    uuid.UUID     `json:"team_id" db:"team_id"`
	UserID    uuid.UUID     `json:"user_id" db:"user_id"`
	AddedBy   uuid.NullUUID `json:"added_by,omitempty" db:"added_by"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	User      *User         `json:"user,omitempty" db:"-"`
}

// Project represents a project
type Project struct {
	Base
	OrganizationID uuid.UUID      `json:"organization_id" db:"organization_id"`
	TeamID         uuid.NullUUID  `json:"team_id,omitempty" db:"team_id"`
	Name           string         `json:"name" db:"name"`
	Description    sql.NullString `json:"description,omitempty" db:"description"`
	Key            sql.NullString `json:"key,omitempty" db:"key"`
	Color          sql.NullString `json:"color,omitempty" db:"color"`
	Icon           sql.NullString `json:"icon,omitempty" db:"icon"`
	Visibility     string         `json:"visibility" db:"visibility"`
	Status         string         `json:"status" db:"status"`
	StartDate      sql.NullTime   `json:"start_date,omitempty" db:"start_date"`
	DueDate        sql.NullTime   `json:"due_date,omitempty" db:"due_date"`
	Settings       []byte         `json:"settings" db:"settings"`
	CreatedBy      uuid.UUID      `json:"created_by" db:"created_by"`
	ArchivedAt     sql.NullTime   `json:"archived_at,omitempty" db:"archived_at"`
}

// ProjectMember represents project membership
type ProjectMember struct {
	Base
	ProjectID uuid.UUID     `json:"project_id" db:"project_id"`
	UserID    uuid.UUID     `json:"user_id" db:"user_id"`
	Role      string        `json:"role" db:"role"`
	AddedBy   uuid.NullUUID `json:"added_by,omitempty" db:"added_by"`
	User      *User         `json:"user,omitempty" db:"-"`
}

// TaskStatus represents a task status
type TaskStatus struct {
	Base
	ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Color       string    `json:"color" db:"color"`
	Position    int       `json:"position" db:"position"`
	IsDefault   bool      `json:"is_default" db:"is_default"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
}

// Task represents a task/issue
type Task struct {
	Base
	ProjectID      uuid.UUID      `json:"project_id" db:"project_id"`
	ParentTaskID   uuid.NullUUID  `json:"parent_task_id,omitempty" db:"parent_task_id"`
	Title          string         `json:"title" db:"title"`
	Description    sql.NullString `json:"description,omitempty" db:"description"`
	TaskNumber     int            `json:"task_number" db:"task_number"`
	StatusID       uuid.UUID      `json:"status_id" db:"status_id"`
	Priority       string         `json:"priority" db:"priority"`
	AssigneeID     uuid.NullUUID  `json:"assignee_id,omitempty" db:"assignee_id"`
	ReporterID     uuid.UUID      `json:"reporter_id" db:"reporter_id"`
	StartDate      sql.NullTime   `json:"start_date,omitempty" db:"start_date"`
	DueDate        sql.NullTime   `json:"due_date,omitempty" db:"due_date"`
	CompletedAt    sql.NullTime   `json:"completed_at,omitempty" db:"completed_at"`
	EstimatedHours sql.NullFloat64 `json:"estimated_hours,omitempty" db:"estimated_hours"`
	ActualHours    float64        `json:"actual_hours" db:"actual_hours"`
	Position       float64        `json:"position" db:"position"`
	CustomFields   []byte         `json:"custom_fields" db:"custom_fields"`
	CreatedBy      uuid.UUID      `json:"created_by" db:"created_by"`

	// Relations (not in DB, loaded separately)
	Status   *TaskStatus `json:"status,omitempty" db:"-"`
	Assignee *User       `json:"assignee,omitempty" db:"-"`
	Reporter *User       `json:"reporter,omitempty" db:"-"`
	Labels   []Label     `json:"labels,omitempty" db:"-"`
	Subtasks []Task      `json:"subtasks,omitempty" db:"-"`
}

// TaskAssignee represents multiple assignees for a task
type TaskAssignee struct {
	ID         uuid.UUID     `json:"id" db:"id"`
	TaskID     uuid.UUID     `json:"task_id" db:"task_id"`
	UserID     uuid.UUID     `json:"user_id" db:"user_id"`
	AssignedBy uuid.NullUUID `json:"assigned_by,omitempty" db:"assigned_by"`
	CreatedAt  time.Time     `json:"created_at" db:"created_at"`
	User       *User         `json:"user,omitempty" db:"-"`
}

// TaskDependency represents task dependencies
type TaskDependency struct {
	ID               uuid.UUID `json:"id" db:"id"`
	TaskID           uuid.UUID `json:"task_id" db:"task_id"`
	DependsOnTaskID  uuid.UUID `json:"depends_on_task_id" db:"depends_on_task_id"`
	DependencyType   string    `json:"dependency_type" db:"dependency_type"`
	CreatedBy        uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// Label represents a label/tag
type Label struct {
	Base
	OrganizationID uuid.UUID      `json:"organization_id" db:"organization_id"`
	ProjectID      uuid.NullUUID  `json:"project_id,omitempty" db:"project_id"`
	Name           string         `json:"name" db:"name"`
	Color          string         `json:"color" db:"color"`
	Description    sql.NullString `json:"description,omitempty" db:"description"`
	CreatedBy      uuid.UUID      `json:"created_by" db:"created_by"`
}

// TaskLabel represents the many-to-many relationship
type TaskLabel struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TaskID    uuid.UUID `json:"task_id" db:"task_id"`
	LabelID   uuid.UUID `json:"label_id" db:"label_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Label     *Label    `json:"label,omitempty" db:"-"`
}

// Comment represents a task comment
type Comment struct {
	Base
	TaskID          uuid.UUID      `json:"task_id" db:"task_id"`
	UserID          uuid.UUID      `json:"user_id" db:"user_id"`
	Content         string         `json:"content" db:"content"`
	ParentCommentID uuid.NullUUID  `json:"parent_comment_id,omitempty" db:"parent_comment_id"`
	IsEdited        bool           `json:"is_edited" db:"is_edited"`
	User            *User          `json:"user,omitempty" db:"-"`
	Replies         []Comment      `json:"replies,omitempty" db:"-"`
}

// ActivityLog represents system activity
type ActivityLog struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	OrganizationID uuid.UUID     `json:"organization_id" db:"organization_id"`
	ProjectID      uuid.NullUUID `json:"project_id,omitempty" db:"project_id"`
	TaskID         uuid.NullUUID `json:"task_id,omitempty" db:"task_id"`
	UserID         uuid.UUID     `json:"user_id" db:"user_id"`
	ActivityType   string        `json:"activity_type" db:"activity_type"`
	Metadata       []byte        `json:"metadata" db:"metadata"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	User           *User         `json:"user,omitempty" db:"-"`
}

// Attachment represents a file attachment
type Attachment struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	TaskID        uuid.NullUUID  `json:"task_id,omitempty" db:"task_id"`
	CommentID     uuid.NullUUID  `json:"comment_id,omitempty" db:"comment_id"`
	UploadedBy    uuid.UUID      `json:"uploaded_by" db:"uploaded_by"`
	FileName      string         `json:"file_name" db:"file_name"`
	FileSize      int64          `json:"file_size" db:"file_size"`
	MimeType      string         `json:"mime_type" db:"mime_type"`
	StoragePath   string         `json:"storage_path" db:"storage_path"`
	ThumbnailPath sql.NullString `json:"thumbnail_path,omitempty" db:"thumbnail_path"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	DeletedAt     sql.NullTime   `json:"deleted_at,omitempty" db:"deleted_at"`
	Uploader      *User          `json:"uploader,omitempty" db:"-"`
}

// TimeEntry represents time tracking
type TimeEntry struct {
	Base
	TaskID      uuid.UUID       `json:"task_id" db:"task_id"`
	UserID      uuid.UUID       `json:"user_id" db:"user_id"`
	Description sql.NullString  `json:"description,omitempty" db:"description"`
	Hours       float64         `json:"hours" db:"hours"`
	StartedAt   time.Time       `json:"started_at" db:"started_at"`
	EndedAt     sql.NullTime    `json:"ended_at,omitempty" db:"ended_at"`
	IsBillable  bool            `json:"is_billable" db:"is_billable"`
	User        *User           `json:"user,omitempty" db:"-"`
}

// CustomField represents custom fields configuration
type CustomField struct {
	Base
	ProjectID    uuid.UUID       `json:"project_id" db:"project_id"`
	Name         string          `json:"name" db:"name"`
	FieldType    string          `json:"field_type" db:"field_type"`
	Description  sql.NullString  `json:"description,omitempty" db:"description"`
	IsRequired   bool            `json:"is_required" db:"is_required"`
	Options      []byte          `json:"options,omitempty" db:"options"`
	DefaultValue []byte          `json:"default_value,omitempty" db:"default_value"`
	Position     int             `json:"position" db:"position"`
	CreatedBy    uuid.UUID       `json:"created_by" db:"created_by"`
}

// View represents saved views/filters
type View struct {
	Base
	ProjectID     uuid.UUID      `json:"project_id" db:"project_id"`
	CreatedByID   uuid.UUID      `json:"created_by" db:"created_by"`
	Name          string         `json:"name" db:"name"`
	ViewType      string         `json:"view_type" db:"view_type"`
	IsDefault     bool           `json:"is_default" db:"is_default"`
	IsShared      bool           `json:"is_shared" db:"is_shared"`
	Filters       []byte         `json:"filters" db:"filters"`
	SortBy        []byte         `json:"sort_by" db:"sort_by"`
	GroupBy       sql.NullString `json:"group_by,omitempty" db:"group_by"`
	VisibleFields []byte         `json:"visible_fields" db:"visible_fields"`
	Settings      []byte         `json:"settings" db:"settings"`
	Position      int            `json:"position" db:"position"`
}

// Notification represents a user notification
type Notification struct {
	ID               uuid.UUID     `json:"id" db:"id"`
	UserID           uuid.UUID     `json:"user_id" db:"user_id"`
	NotificationType string        `json:"notification_type" db:"notification_type"`
	Title            string        `json:"title" db:"title"`
	Message          sql.NullString `json:"message,omitempty" db:"message"`
	TaskID           uuid.NullUUID `json:"task_id,omitempty" db:"task_id"`
	ProjectID        uuid.NullUUID `json:"project_id,omitempty" db:"project_id"`
	TriggeredBy      uuid.NullUUID `json:"triggered_by,omitempty" db:"triggered_by"`
	IsRead           bool          `json:"is_read" db:"is_read"`
	ReadAt           sql.NullTime  `json:"read_at,omitempty" db:"read_at"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	Task             *Task         `json:"task,omitempty" db:"-"`
	Project          *Project      `json:"project,omitempty" db:"-"`
	Triggerer        *User         `json:"triggerer,omitempty" db:"-"`
}
