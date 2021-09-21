package workspace

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

// CreateAsset adds new Asset entity to the database.
// If creation is successful, the method returns Asset entity.
// Can return validation or database errors.
func (str Store) CreateAsset(ctx context.Context, claims auth.Claims, newAsset NewAsset, now time.Time) (Asset,
	error) {
	if err := validation.Check(ctx, newAsset); err != nil {
		return Asset{}, fmt.Errorf("error during data validation of Asset entity: %w", err)
	}

	asset := Asset{
		ID:            uuid.Generate(),
		WorkspaceID:   newAsset.WorkspaceID,
		AssetRefID:    newAsset.AssetRefID,
		X:             newAsset.X,
		Y:             newAsset.Y,
		Z:             newAsset.Z,
		Scale:         newAsset.Scale,
		Height:        newAsset.Height,
		Width:         newAsset.Width,
		Length:        newAsset.Length,
		DateCreated:   now,
		DateUpdated:   now,
		CreatedByUser: claims.Subject,
		UpdatedByUser: claims.Subject,
	}

	const query = `
	INSERT INTO ASSET
		(asset_id, workspace_id, asset_external_ref_id, position_x, position_y, position_z, scale, height_by_y,
			width_by_x, length_by_z, date_created, created_by_user_id, date_updated, updated_by_user_id)
	VALUES
		(:asset_id, :workspace_id, :asset_external_ref_id, :position_x, :position_y, :position_z, :scale, :height_by_y,
			:width_by_x, :length_by_z, :date_created, :created_by_user_id, :date_updated, :updated_by_user_id)`

	if err := database.NamedExecContext(ctx, str.logger, str.connection, query, asset); err != nil {
		return Asset{}, fmt.Errorf("error during create of new Asset entity: %w", err)
	}

	return asset, nil
}

// UpdateAsset change existing Asset entity in the database.
// If error occurs, the method can return validation or database errors.
func (str Store) UpdateAsset(ctx context.Context, claims auth.Claims, assetId string, asset UpdateAsset,
	now time.Time) error {
	if err := uuid.Validate(assetId); err != nil {
		return database.ErrorInvalidIdentifier
	}

	if err := validation.Check(ctx, asset); err != nil {
		return fmt.Errorf("error during data validation of Asset entity: %w", err)
	}

	assetData, err := str.QueryAssetByID(ctx, assetId)
	if err != nil {
		return fmt.Errorf("error during search of Asset entity -> id={%q}: %w", assetId, err)
	}

	if assetData.CreatedByUser != claims.Subject {
		return database.ErrorForbidden
	}

	return errors.New("TODO IMPLEMENT RECEIVER")
}

// DeleteAsset removes existing Asset entity in the database.
// If error occurs, the method can return database errors.
func (str Store) DeleteAsset(ctx context.Context, claims auth.Claims, assetId string) error {
	if err := uuid.Validate(assetId); err != nil {
		return database.ErrorInvalidIdentifier
	}

	assetData, err := str.QueryAssetByID(ctx, assetId)
	if err != nil {
		return fmt.Errorf("error during search of Asset entity -> id={%q}: %w", assetId, err)
	}

	if assetData.CreatedByUser != claims.Subject {
		return database.ErrorForbidden
	}

	queryParams := struct {
		AssetID string `db:"asset_id"`
	}{
		AssetID: assetId,
	}

	const query = `
	DELETE FROM
		ASSET
	WHERE
		asset_id = :asset_id`

	if err := database.NamedExecContext(ctx, str.logger, str.connection, query, queryParams); err != nil {
		return fmt.Errorf("error during delete of Asset entity -> id={%q}: %w", assetId, err)
	}

	return nil
}

// QueryWorkspaceAssets looking for all Asset entities that belong to a specific Workspace using skip/top mechanics with
// descending order by update date field.
func (str Store) QueryWorkspaceAssets(ctx context.Context, workspaceId string, skip int32, top int32) ([]Asset,
	error) {
	if err := uuid.Validate(workspaceId); err != nil {
		return []Asset{}, database.ErrorInvalidIdentifier
	}

	queryParams := struct {
		WorkspaceId string `db:"workspace_id"`
		Skip        int32  `db:"offset"`
		Top         int32  `db:"top"`
	}{
		WorkspaceId: workspaceId,
		Skip:        skip,
		Top:         top,
	}

	const query = `
	SELECT
		a.asset_id,
		a.workspace_id,
		a.asset_external_ref_id,
		a.position_x,
		a.position_y,
		a.position_z,
		a.scale,
		a.height_by_y,
		a.width_by_x,
		a.length_by_z,
		a.date_created,
		a.created_by_user_id,
		a.date_updated,
		a.updated_by_user_id
	FROM
		ASSET AS a
	WHERE
		a.workspace_id = :workspace_id
	ORDER BY a.date_updated DESC
	OFFSET :offset ROWS FETCH NEXT :top ROWS ONLY`

	var assetCollection []Asset
	if err := database.NamedQuerySlice(ctx, str.logger, str.connection, query, queryParams, &assetCollection); err != nil {
		if err == database.ErrorNotFound {
			return nil, database.ErrorNotFound
		}

		return nil, fmt.Errorf("error during search of Asset entities: %w", err)
	}

	return assetCollection, nil
}

// QueryAssetByID looking for Asset entity with assetId identifier.
func (str Store) QueryAssetByID(ctx context.Context, assetId string) (Asset, error) {
	if err := uuid.Validate(assetId); err != nil {
		return Asset{}, err
	}

	queryParams := struct {
		AssetID string `db:"asset_id"`
	}{
		AssetID: assetId,
	}

	const query = `
	SELECT
		a.asset_id,
		a.workspace_id,
		a.asset_external_ref_id,
		a.position_x,
		a.position_y,
		a.position_z,
		a.scale,
		a.height_by_y,
		a.width_by_x,
		a.length_by_z,
		a.date_created,
		a.created_by_user_id,
		a.date_updated,
		a.updated_by_user_id
	FROM
		ASSET AS a
	WHERE
		p.asset_id = :asset_id`

	var assetData Asset
	if err := database.NamedQueryStruct(ctx, str.logger, str.connection, query, queryParams, &assetData); err != nil {
		if err == database.ErrorNotFound {
			return Asset{}, database.ErrorNotFound
		}

		return Asset{}, fmt.Errorf("error during search of Asset entity -> id={%q}: %w", assetId, err)
	}

	return assetData, nil
}
