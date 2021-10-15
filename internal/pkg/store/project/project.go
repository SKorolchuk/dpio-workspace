package project

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/auth"
	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/database"
	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/uuid"
	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/validation"
)

// CreateProject adds new Project entity to the database.
// If creation is successful, the method returns Project entity.
// Can return validation or database errors.
func (str Store) CreateProject(ctx context.Context, claims auth.Claims, project NewProject, now time.Time) (Project,
	error) {
	if err := validation.Check(ctx, project); err != nil {
		return Project{}, fmt.Errorf("error during data validation of Project entity: %w", err)
	}

	projectData := Project{
		ID:            uuid.Generate(),
		ProjectTypeID: project.ProjectTypeID,
		Name:          project.Name,
		Description:   project.Description,
		DateCreated:   now,
		DateUpdated:   now,
		CreatedByUser: claims.Subject,
		UpdatedByUser: claims.Subject,
	}

	const query = `
	INSERT INTO PROJECT
		(project_id, project_collaboration_type_id, name, description, date_created,
			created_by_user_id, date_updated, updated_by_user_id)
	VALUES
		(:project_id, :project_collaboration_type_id, :name, :description, :date_created,
			:created_by_user_id, :date_updated, :updated_by_user_id)`

	if err := database.NamedExecContext(ctx, str.logger, str.connection, query, projectData); err != nil {
		return Project{}, fmt.Errorf("error during create of new Project entity: %w", err)
	}

	return projectData, nil
}

// UpdateProject change existing Project entity in the database.
// If error occurs, the method can return validation or database errors.
func (str Store) UpdateProject(ctx context.Context, claims auth.Claims, projectId string, project UpdateProject,
	now time.Time) error {
	if err := uuid.Validate(projectId); err != nil {
		return database.ErrorInvalidIdentifier
	}

	if err := validation.Check(ctx, project); err != nil {
		return fmt.Errorf("error during data validation of Project entity: %w", err)
	}

	projectData, err := str.QueryProjectByID(ctx, projectId)
	if err != nil {
		return fmt.Errorf("error during search of Project entity -> id={%q}: %w", projectId, err)
	}

	// TODO implement group-based edit check in downstream logic
	if projectData.CreatedByUser != claims.Subject {
		return database.ErrorForbidden
	}

	if project.Name != nil {
		projectData.Name = *project.Name
	}

	if project.Description != nil {
		projectData.Description = *project.Description
	}

	projectData.DateUpdated = now
	projectData.UpdatedByUser = claims.Subject

	const query = `
	UPDATE
		PROJECT
	SET
		"name" = :name,
		"description" = :description,
		"date_updated" = :date_updated,
		"updated_by_user_id" = :updated_by_user_id
	WHERE
		project_id = :project_id`

	if err := database.NamedExecContext(ctx, str.logger, str.connection, query, projectData); err != nil {
		return fmt.Errorf("error during update of Project entity -> id={%s}: %w", projectId, err)
	}

	return nil
}

// DeleteProject removes existing Project entity in the database.
// If error occurs, the method can return database errors.
func (str Store) DeleteProject(ctx context.Context, claims auth.Claims, projectId string) error {
	if err := uuid.Validate(wsId); err != nil {
		return database.ErrorInvalidIdentifier
	}

	wsData, err := str.QueryWorkspaceByID(ctx, wsId)
	if err != nil {
		return fmt.Errorf("error during search of Workspace entity -> id={%q}: %w", wsId, err)
	}

	if wsData.CreatedByUser != claims.Subject {
		return database.ErrorForbidden
	}

	queryParams := struct {
		WorkspaceID string `db:"workspace_id"`
	}{
		WorkspaceID: wsId,
	}

	// TODO Implement cascade delete of Assets and return list of marked Asset External Refs

	const query = `
	DELETE FROM
		WORKSPACE
	WHERE
		workspace_id = :workspace_id`

	if err := database.NamedExecContext(ctx, str.logger, str.connection, query, queryParams); err != nil {
		return fmt.Errorf("error during delete of Workspace entity -> id={%q}: %w", wsId, err)
	}

	return nil
}

// QueryProjects looking for all Project entities using skip/top mechanics with descending order by update date field.
func (str Store) QueryProjects(ctx context.Context, skip int32, top int32) ([]Project, error) {
	return []Project{}, errors.New("TODO Implement method")
}

// QueryProjectByID looking for Project entity with projectId identifier.
func (str Store) QueryProjectByID(ctx context.Context, projectId string) (Project, error) {
	if err := uuid.Validate(projectId); err != nil {
		return Project{}, err
	}

	queryParams := struct {
		ProjectID string `db:"project_id"`
	}{
		ProjectID: projectId,
	}

	const query = `
	SELECT
		p.project_id,
		p.project_collaboration_type_id,
		p.name,
		p.description,
		p.date_created,
		p.created_by_user_id,
		p.date_updated,
		p.updated_by_user_id
	FROM
		PROJECT AS p
	WHERE
		p.project_id = :project_id`

	var projectData Project
	if err := database.NamedQueryStruct(ctx, str.logger, str.connection, query, queryParams, &projectData); err != nil {
		if err == database.ErrorNotFound {
			return Project{}, database.ErrorNotFound
		}

		return Project{}, fmt.Errorf("error during search of Project entity -> id={%q}: %w", projectId, err)
	}

	return projectData, nil
}
