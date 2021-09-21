package workspace

import "time"

const (
	StickerWorkspaceType       = "2fdf996e-2372-4f3c-bccf-d8efcca8bd49"
	Environment2DWorkspaceType = "78e95523-4ed2-49e6-8b1a-b8c073daab41"
	Environment3DWorkspaceType = "b8d78dda-027c-498e-8609-33cc6f4a6dbe"
)

// Stem represents specific type of Workspace entities.
type Stem struct {
	ID   string `db:"stem_id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Workspace represents shared area when Asset objects can be composed into a scene or a part of specific project.
type Workspace struct {
	ID               string    `db:"workspace_id" json:"id"`
	ProjectID        string    `db:"project_id" json:"projectId"`
	StemID           string    `db:"stem_id" json:"stemId"`
	Name             string    `db:"name" json:"name"`
	Description      string    `db:"description" json:"description"`
	AssetAmountLimit int32     `db:"asset_amount_limit" json:"assetAmountLimit"`
	MaxX             int32     `db:"x_max" json:"maxX"`
	MaxY             int32     `db:"y_max" json:"maxY"`
	MaxZ             int32     `db:"z_max" json:"maxZ"`
	DateCreated      time.Time `db:"date_created" json:"dateCreated"`
	CreatedByUser    string    `db:"created_by_user_id" json:"createdByUser"`
	DateUpdated      time.Time `db:"date_updated" json:"dateUpdated"`
	UpdatedByUser    string    `db:"updated_by_user_id" json:"updatedByUser"`
}

// NewWorkspace describes all data that should be specified during creation of new Workspace entity.
type NewWorkspace struct {
	ProjectID        string `json:"projectId" validate:"required"`
	StemID           string `json:"stemId" validate:"required"`
	Name             string `json:"name" validate:"required"`
	Description      string `json:"description" validate:"-"`
	AssetAmountLimit int32  `json:"assetAmountLimit" validate:"required,gte=1"`
	MaxX             int32  `json:"maxX" validate:"required,gte=1"`
	MaxY             int32  `json:"maxY" validate:"required,gte=1"`
	MaxZ             int32  `json:"maxZ" validate:"required,gte=0"`
}

// UpdateWorkspace describes all data that can be changed during update of existing Workspace entity.
type UpdateWorkspace struct {
	Name             *string `json:"name" validate:"required"`
	Description      *string `json:"description" validate:"-"`
	AssetAmountLimit *int32  `json:"assetAmountLimit" validate:"required,gte=1"`
	MaxX             *int32  `json:"maxX" validate:"required,gte=1"`
	MaxY             *int32  `json:"maxY" validate:"required,gte=1"`
	MaxZ             *int32  `json:"maxZ" validate:"required,gte=0"`
}

// Asset represents reference to dynamic element that should be displayed in Workspace area.
type Asset struct {
	ID            string    `db:"asset_id" json:"id"`
	WorkspaceID   string    `db:"workspace_id" json:"workspaceId"`
	AssetRefID    string    `db:"asset_external_ref_id" json:"assetRefId"`
	X             int32     `db:"position_x" json:"x"`
	Y             int32     `db:"position_y" json:"y"`
	Z             int32     `db:"position_z" json:"z"`
	Scale         int32     `db:"scale" json:"scale"`
	Height        int32     `db:"height_by_y" json:"height"`
	Width         int32     `db:"width_by_x" json:"width"`
	Length        int32     `db:"length_by_z" json:"length"`
	DateCreated   time.Time `db:"date_created" json:"dateCreated"`
	CreatedByUser string    `db:"created_by_user_id" json:"createdByUser"`
	DateUpdated   time.Time `db:"date_updated" json:"dateUpdated"`
	UpdatedByUser string    `db:"updated_by_user_id" json:"updatedByUser"`
}

// NewAsset describes all data that should be specified during creation of new Asset entity.
type NewAsset struct {
	WorkspaceID string `json:"workspaceId" validate:"required"`
	AssetRefID  string `json:"assetRefId" validate:"required"`
	X           int32  `json:"x" validate:"required,gte=0"`
	Y           int32  `json:"y" validate:"required,gte=0"`
	Z           int32  `json:"z" validate:"required,gte=0"`
	Scale       int32  `json:"scale" validate:"required,gte=0"`
	Height      int32  `json:"height" validate:"required,gte=0"`
	Width       int32  `json:"width" validate:"required,gte=0"`
	Length      int32  `json:"length" validate:"required,gte=0"`
}

// UpdateAsset describes all data that can be changed during update of existing Asset entity.
type UpdateAsset struct {
	// TODO Implement structure
}
