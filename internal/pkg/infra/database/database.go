package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/server"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // init function.
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	DefaultRetryPeriodInMilliseconds = 100 * time.Millisecond
)

// Common errors for CRUD operations.
var (
	ErrorNotFound          = errors.New("entity not found")
	ErrorInvalidIdentifier = errors.New("specified ID format is not valid")
	ErrorAuthFail          = errors.New("authentication error")
	ErrorForbidden         = errors.New("action is not allowed")
)

// DbConfig is used by PostgreSQL driver to build connection string and establish database connection.
type DbConfig struct {
	User               string
	Password           string
	Host               string
	DatabaseName       string
	MaxIdleConnections int32
	MaxOpenConnections int32
	DisableTLS         bool
}

// Open use DbConfig settings to open database connection.
func Open(config DbConfig) (*sqlx.DB, error) {
	if err := validateDbConfig(config); err != nil {
		return nil, err
	}

	sslMode := "require"

	if config.DisableTLS {
		sslMode = "disable"
	}

	connectionParameters := make(url.Values)
	connectionParameters.Set("sslmode", sslMode)
	connectionParameters.Set("timezone", "utc")

	connectionString := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(config.User, config.Password),
		Host:     config.Host,
		Path:     config.DatabaseName,
		RawQuery: connectionParameters.Encode(),
	}

	dbConnection, err := sqlx.Open("postgres", connectionString.String())
	if err != nil {
		return nil, err
	}

	dbConnection.SetMaxIdleConns(int(config.MaxIdleConnections))
	dbConnection.SetMaxOpenConns(int(config.MaxOpenConnections))

	return dbConnection, nil
}

// StatusCheck returns error if issues exist with database connection.
func StatusCheck(ctx context.Context, connection *sqlx.DB) error {
	var connectivityError error

	for attempt := 1; ; attempt++ {
		connectivityError = connection.PingContext(ctx)
		if connectivityError == nil {
			break
		}
		time.Sleep(time.Duration(attempt) * DefaultRetryPeriodInMilliseconds)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	const pingQuery = `SELECT true`
	var output bool
	return connection.QueryRowContext(ctx, pingQuery).Scan(&output)
}

// NamedExecContext is a helper to execute a CRUD operation under logging and tracing features.
func NamedExecContext(ctx context.Context, logger *zap.SugaredLogger, connection *sqlx.DB, sqlQuery string,
	params interface{}) error {
	query, err := queryString(sqlQuery, params)
	if err != nil {
		return err
	}
	logger.Infow("database.NameExecContext", "traceid", server.GetTraceID(ctx), "query", query)

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "database.query")
	span.SetAttributes(attribute.String("query", query))
	defer span.End()

	if _, err := connection.NamedExecContext(ctx, query, params); err != nil {
		return err
	}

	return nil
}

// NamedQueryStruct is a helper to execute queries that return a single structured value.
func NamedQueryStruct(ctx context.Context, logger *zap.SugaredLogger, connection *sqlx.DB, sqlQuery string,
	params interface{}, target interface{}) error {
	query, err := queryString(sqlQuery, params)
	if err != nil {
		return err
	}
	logger.Infow("database.NamedQueryStruct", "traceid", server.GetTraceID(ctx), "query", query)

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "database.query")
	span.SetAttributes(attribute.String("query", query))
	defer span.End()

	rows, err := connection.NamedQueryContext(ctx, query, params)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return ErrorNotFound
	}

	if err := rows.StructScan(target); err != nil {
		return err
	}

	return nil
}

// NamedQuerySlice is a helper to execute queries that return a collection of data.
func NamedQuerySlice(ctx context.Context, logger *zap.SugaredLogger, connection *sqlx.DB, sqlQuery string,
	params interface{}, target interface{}) error {
	query, err := queryString(sqlQuery, params)
	if err != nil {
		return err
	}
	logger.Infow("database.NamedQuerySlice", "traceid", server.GetTraceID(ctx), "query", query)

	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "database.query")
	span.SetAttributes(attribute.String("query", query))
	defer span.End()

	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Slice {
		return errors.New("target object should be a pointer to a slice")
	}

	rows, err := connection.NamedQueryContext(ctx, query, params)
	if err != nil {
		return err
	}

	sliceRef := value.Elem()
	for rows.Next() {
		sliceElement := reflect.New(sliceRef.Type().Elem())
		if err := rows.StructScan(sliceElement.Interface()); err != nil {
			return err
		}
		sliceRef.Set(reflect.Append(sliceRef, sliceElement.Elem()))
	}

	return nil
}

// queryString formats SQL query and set specified arguments.
func queryString(sqlQuery string, args ...interface{}) (string, error) {
	query, params, err := sqlx.Named(sqlQuery, args)

	if err != nil {
		return "", err
	}

	for _, param := range params {
		var value string

		switch valueType := param.(type) {
		case string:
			value = fmt.Sprintf("%q", valueType)
		case []byte:
			value = fmt.Sprintf("%q", string(valueType))
		default:
			value = fmt.Sprintf("%v", valueType)
		}

		query = strings.Replace(query, "?", value, 1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.Trim(query, " "), nil
}

func validateDbConfig(config DbConfig) error {
	if len(config.User) <= 0 {
		return errors.New("username is not specified in connection settings")
	}

	if len(config.Password) <= 0 {
		return errors.New("password is not specified in connection settings")
	}

	if len(config.Host) <= 0 {
		return errors.New("database host is not specified in connection settings")
	}

	if len(config.DatabaseName) <= 0 {
		return errors.New("database name is not specified in connection settings")
	}

	return nil
}
