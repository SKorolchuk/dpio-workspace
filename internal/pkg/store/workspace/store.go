package workspace

import (
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
