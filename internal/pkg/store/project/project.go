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
	return errors.New("TODO Implement method")
}

// DeleteProject removes existing Project entity in the database.
// If error occurs, the method can return database errors.
func (str Store) DeleteProject(ctx context.Context, claims auth.Claims, projectId string) error {
	return errors.New("TODO Implement method")
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
