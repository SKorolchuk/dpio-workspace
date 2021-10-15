package database

import (
	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/auth"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type TestDBConfig struct {
	ContainerImageName string
	ExposedPort string
	AdditionalArguments []string
}

type TestDBConnection struct {
	Connection *sqlx.DB
	Logger *zap.SugaredLogger
	Session *auth.Auth
}
