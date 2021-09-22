package project

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store represents a point of access to CollaborationType, Project, Role, Group, GroupRole, GroupUser and
// GroupAccess entities.
type Store struct {
	logger     *zap.SugaredLogger
	connection *sqlx.DB
}

// NewStore creates an instance of Store for access to CollaborationType, Project, Role, Group, GroupRole, GroupUser and
// GroupAccess entities.
func NewStore(logger *zap.SugaredLogger, connection *sqlx.DB) Store {
	return Store{
		logger:     logger,
		connection: connection,
	}
}
