package project

import "time"

// CollaborationType represents a team collaboration level for a Project entity.
type CollaborationType struct {
	ID   string `db:"project_collaboration_type_id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Project represents a unit of work structuring.
type Project struct {
	ID            string    `db:"project_id" json:"id"`
	ProjectTypeID string    `db:"project_collaboration_type_id" json:"projectTypeId"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	DateCreated   time.Time `db:"date_created" json:"dateCreated"`
	CreatedByUser string    `db:"created_by_user_id" json:"createdByUser"`
	DateUpdated   time.Time `db:"date_updated" json:"dateUpdated"`
	UpdatedByUser string    `db:"updated_by_user_id" json:"updatedByUser"`
}

// Role represents a one or several permissions to allow/restrict work in Project.
type Role struct {
	ID   string `db:"project_role_id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Group represents a unit of team organization.
type Group struct {
	ID            string    `db:"project_group_id" json:"id"`
	Name          string    `db:"name" json:"name"`
	DateCreated   time.Time `db:"date_created" json:"dateCreated"`
	CreatedByUser string    `db:"created_by_user_id" json:"createdByUser"`
	DateUpdated   time.Time `db:"date_updated" json:"dateUpdated"`
	UpdatedByUser string    `db:"updated_by_user_id" json:"updatedByUser"`
}

// GroupRole represents a Role assignment to a Group.
type GroupRole struct {
	ID            string    `db:"project_group_role_id" json:"id"`
	GroupID       string    `db:"project_group_id" json:"groupId"`
	RoleID        string    `db:"project_role_id" json:"roleId"`
	DateCreated   time.Time `db:"date_created" json:"dateCreated"`
	CreatedByUser string    `db:"created_by_user_id" json:"createdByUser"`
}

// GroupUser represents a user assignment to a Group.
type GroupUser struct {
	ID            string    `db:"project_group_user_id" json:"id"`
	GroupID       string    `db:"project_group_id" json:"groupId"`
	UserID        string    `db:"user_id" json:"userId"`
	DateCreated   time.Time `db:"date_created" json:"dateCreated"`
	CreatedByUser string    `db:"created_by_user_id" json:"createdByUser"`
}

// GroupAccess represents a Group assignment to a Project.
type GroupAccess struct {
	ID            string    `db:"project_group_access_id" json:"id"`
	ProjectID     string    `db:"project_id" json:"projectId"`
	GroupID       string    `db:"project_group_id" json:"groupId"`
	DateCreated   time.Time `db:"date_created" json:"dateCreated"`
	CreatedByUser string    `db:"created_by_user_id" json:"createdByUser"`
}

// NewProject describes all data that should be specified during creation of new Project entity.
type NewProject struct {
	ProjectTypeID string `json:"projectTypeId" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description" validate:"required"`
}

// UpdateProject describes all data that can be changed during update of existing Project entity.
type UpdateProject struct {
	// TODO Implement structure
}
