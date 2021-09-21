package workspace

import (
	"context"
	"fmt"

	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/database"
	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/uuid"
)

// QueryStems looking for all Stem entities.
func (str Store) QueryStems(ctx context.Context, skip int32, top int32) ([]Stem, error) {
	queryParams := struct {
		Skip int32 `db:"offset"`
		Top  int32 `db:"top"`
	}{
		Skip: skip,
		Top:  top,
	}

	const query = `
	SELECT
		s.stem_id,
		s.name
	FROM
		STEM AS s
	ORDER BY s.stem_id DESC
	OFFSET :offset ROWS FETCH NEXT :top ROWS ONLY`

	var stemCollection []Stem
	if err := database.NamedQuerySlice(ctx, str.logger, str.connection, query, queryParams, &stemCollection); err != nil {
		if err == database.ErrorNotFound {
			return nil, database.ErrorNotFound
		}

		return nil, fmt.Errorf("error during search of Stem entities: %w", err)
	}

	return stemCollection, nil
}

// QueryStemByID looking for Stem entity with stemId identifier.
func (str Store) QueryStemByID(ctx context.Context, stemId string) (Stem, error) {
	if err := uuid.Validate(stemId); err != nil {
		return Stem{}, err
	}

	queryParams := struct {
		StemID string `db:"stem_id"`
	}{
		StemID: stemId,
	}

	const query = `
	SELECT
		s.stem_id,
		s.name
	FROM
		STEM AS s
	WHERE
		p.stem_id = :stem_id`

	var stem Stem
	if err := database.NamedQueryStruct(ctx, str.logger, str.connection, query, queryParams, &stem); err != nil {
		if err == database.ErrorNotFound {
			return Stem{}, database.ErrorNotFound
		}

		return Stem{}, fmt.Errorf("error during search of Stem entity -> id={%q}: %w", stemId, err)
	}

	return stem, nil
}
