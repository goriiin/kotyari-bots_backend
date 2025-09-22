package api_integrations

import (
	"context"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (a *APIIntegrationsRepo) GetIntegrations(ctx context.Context, integrationName string) ([]model.APIIntegration, error) {
	const query = `
		select provider, url 
		from integrations 
		where provider = $1;
	`

	rows, err := a.db.Query(ctx, query, integrationName)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrapf(err, "no rows for %s", integrationName)
		}

		return nil, errors.Wrap(err, "unexpected error happened")
	}

	apiIntegrations, err := pgx.CollectRows(rows, pgx.RowToStructByName[APIIntegrationDTO])
	if err != nil {
		return nil, errors.Wrap(err, "failed to collect rows")
	}

	if len(apiIntegrations) == 0 {
		return nil, errors.New("no rows for this integration")
	}

	return apiIntegrationToModelSlice(apiIntegrations), nil
}
