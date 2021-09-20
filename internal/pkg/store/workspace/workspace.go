// Package workspace contains CRUD operation related to Workspace core entity.
package workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/auth"
	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/database"
	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/uuid"
	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/validation"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store represents a point of access to Workspace, Asset and Stem entities.
type Store struct {
	logger     *zap.SugaredLogger
	connection *sqlx.DB
}

// NewStore creates an instance of Store for access to Workspace, Asset and Stem entities.
func NewStore(logger *zap.SugaredLogger, connection *sqlx.DB) Store {
	return Store{
		logger:     logger,
		connection: connection,
	}
}

// CreateWorkspace adds new Workspace entity to the database.
// If creation is successful, the method returns Workspace entity.
// Can return validation or database errors.
func (str Store) CreateWorkspace(ctx context.Context, claims auth.Claims, ws NewWorkspace, now time.Time) (Workspace,
	error) {
	if err := validation.Check(ctx, ws); err != nil {
		return Workspace{}, fmt.Errorf("error during data validation of Workspace entity: %w", err)
	}

	wsData := Workspace{
		ID:               uuid.Generate(),
		ProjectID:        ws.ProjectID,
		StemID:           ws.StemID,
		Name:             ws.Name,
		Description:      ws.Description,
		AssetAmountLimit: ws.AssetAmountLimit,
		MaxX:             ws.MaxX,
		MaxY:             ws.MaxY,
		MaxZ:             ws.MaxZ,
		DateCreated:      now,
		DateUpdated:      now,
		CreatedByUser:    claims.Subject,
		UpdatedByUser:    claims.Subject,
	}

	const query = `
	INSERT INTO WORKSPACE
		(workspace_id, project_id, workspace_type_id, name, description, asset_amount_limit, x_max, y_max, z_max,
			date_created, created_by_user_id, date_updated, updated_by_user_id)
	VALUES (:workspace_id, :project_id, :workspace_type_id, :name, :description, :asset_amount_limit, :x_max, :y_max,
				:z_max, :date_created, :created_by_user_id, :date_updated, :updated_by_user_id)`

	if err := database.NamedExecContext(ctx, str.logger, str.connection, query, wsData); err != nil {
		return Workspace{}, fmt.Errorf("error during create of new Workspace entity: %w", err)
	}

	return wsData, nil
}

// UpdateWorkspace change existing Workspace entity in the database.
// If error occurs, the method can return validation or database errors.
func (str Store) UpdateWorkspace(ctx context.Context, claims auth.Claims, wsId string, ws UpdateWorkspace,
	now time.Time) error {
	if err := uuid.Validate(wsId); err != nil {
		return database.ErrorInvalidIdentifier
	}

	if err := validation.Check(ctx, ws); err != nil {
		return fmt.Errorf("error during data validation of Workspace entity: %w", err)
	}

	wsData, err := str.QueryWorkspaceByID(ctx, wsId)
	if err != nil {
		return fmt.Errorf("error during search of Workspace entity -> id={%q}: %w", wsId, err)
	}

	if wsData.CreatedByUser != claims.Subject {
		return database.ErrorForbidden
	}

	if ws.Name != nil {
		wsData.Name = *ws.Name
	}

	if ws.Description != nil {
		wsData.Description = *ws.Description
	}

	if ws.AssetAmountLimit != nil {
		wsData.AssetAmountLimit = *ws.AssetAmountLimit
	}

	if ws.MaxX != nil {
		wsData.MaxX = *ws.MaxX
	}
	if ws.MaxY != nil {
		wsData.MaxY = *ws.MaxY
	}
	if ws.MaxZ != nil {
		wsData.MaxZ = *ws.MaxZ
	}
	wsData.DateUpdated = now
	wsData.UpdatedByUser = claims.Subject

	const query = `
	UPDATE
		WORKSPACE
	SET
		"name" = :name,
		"description" = :description,
		"asset_amount_limit" = :asset_amount_limit,
		"x_max" = :x_max,
		"y_max" = :y_max,
		"z_max" = :z_max,
		"date_updated" = :date_updated,
		"updated_by_user_id" = :updated_by_user_id
	WHERE
		workspace_id = :workspace_id`

	if err := database.NamedExecContext(ctx, str.logger, str.connection, query, wsData); err != nil {
		return fmt.Errorf("error during update of Workspace entity -> id={%s}: %w", wsId, err)
	}

	return nil
}

// DeleteWorkspace removes existing Workspace entity in the database.
// If error occurs, the method can return database errors.
func (str Store) DeleteWorkspace(ctx context.Context, claims auth.Claims, wsId string) error {
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

// QueryWorkspaces looking for all Workspace entities.
func (str Store) QueryWorkspaces(ctx context.Context, skip int32, top int32) ([]Workspace, error) {
	queryParams := struct {
		Skip int32 `db:"offset"`
		Top  int32 `db:"top"`
	}{
		Skip: skip,
		Top:  top,
	}

	const query = `
	SELECT
		w.workspace_id,
		w.project_id,
		w.workspace_type_id,
		w.name,
		w.description,
		w.asset_amount_limit,
		w.x_max,
		w.y_max,
		w.z_max,
		w.date_created,
		w.created_by_user_id,
		w.date_updated,
		w.updated_by_user_id
	FROM
		WORKSPACE AS w
	OFFSET :offset ROWS FETCH NEXT :top ROWS ONLY`

	var wsCollection []Workspace
	if err := database.NamedQuerySlice(ctx, str.logger, str.connection, query, queryParams, &wsCollection); err != nil {
		if err == database.ErrorNotFound {
			return nil, database.ErrorNotFound
		}

		return nil, fmt.Errorf("error during search of Workspace entities: %w", err)
	}

	return wsCollection, nil
}

// QueryWorkspaceByID looking for Workspace entity with wsId identifier.
func (str Store) QueryWorkspaceByID(ctx context.Context, wsId string) (Workspace, error) {
	if err := uuid.Validate(wsId); err != nil {
		return Workspace{}, err
	}

	queryParams := struct {
		WorkspaceID string `db:"workspace_id"`
	}{
		WorkspaceID: wsId,
	}

	const query = `
	SELECT
		w.workspace_id,
		w.project_id,
		w.workspace_type_id,
		w.name,
		w.description,
		w.asset_amount_limit,
		w.x_max,
		w.y_max,
		w.z_max,
		w.date_created,
		w.created_by_user_id,
		w.date_updated,
		w.updated_by_user_id
	FROM
		WORKSPACE AS w
	WHERE
		p.workspace_id = :workspace_id`

	var wsData Workspace
	if err := database.NamedQueryStruct(ctx, str.logger, str.connection, query, queryParams, &wsData); err != nil {
		if err == database.ErrorNotFound {
			return Workspace{}, database.ErrorNotFound
		}

		return Workspace{}, fmt.Errorf("error during search of Workspace entity -> id={%q}: %w", wsId, err)
	}

	return wsData, nil
}

// TODO add Stem and Asset store operations
