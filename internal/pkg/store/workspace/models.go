package workspace

import "time"

const (
	StickerWorkspaceType       = "2fdf996e-2372-4f3c-bccf-d8efcca8bd49"
	Environment2DWorkspaceType = "78e95523-4ed2-49e6-8b1a-b8c073daab41"
	Environment3DWorkspaceType = "b8d78dda-027c-498e-8609-33cc6f4a6dbe"
)

// Stem represents specific type of Workspace entities.
type Stem struct {
	ID   string `db:"workspace_type_id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Workspace represents shared area when Asset objects can be composed into a scene or a part of specific project.
type Workspace struct {
	ID               string    `db:"workspace_id" json:"id"`
	ProjectID        string    `db:"project_id" json:"projectId"`
	StemID           string    `db:"workspace_type_id" json:"stemId"`
	Name             string    `db:"name" json:"name"`
	Description      string    `db:"description" json:"description"`
	AssetAmountLimit int32     `db:"asset_amount_limit" json:"assetAmountLimit"`
	MaxX             int32     `db:"x_max" json:"MaxX"`
	MaxY             int32     `db:"y_max" json:"MaxY"`
	MaxZ             int32     `db:"z_max" json:"MaxZ"`
	DateCreated      time.Time `db:"date_created" json:"dateCreated"`
	CreatedByUser    string    `db:"created_by_user_id" json:"createdByUser"`
	DateUpdated      time.Time `db:"date_updated" json:"dateUpdated"`
	UpdatedByUser    string    `db:"updated_by_user_id" json:"updatedByUser"`
}
