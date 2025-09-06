package api_integrations

import (
	"context"
	//"errors"
	//"fmt"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (a *APIIntegrationsRepo) GetIntegrations(ctx context.Context, integrationName string) ([]model.APIIntegration, error) {
	a.log.Debug().Msg("integrations repo")

	const query = `
		select provider, url 
		from integrations 
		where provider = $1;
	`

	rows, err := a.db.Query(ctx, query, integrationName)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// TODO: log no rows
			// TODO: errors
			return nil, errors.Wrapf(err, "no rows for %s", integrationName)

			// return nil, fmt.Errorf("no rows for %s", integrationName)
		}
		// TODO: log err
		return nil, errors.Wrap(err, "unexpected error happened")
	}

	apiIntegrations, err := pgx.CollectRows(rows, pgx.RowToStructByName[APIIntegrationDTO])
	if err != nil {
		// TODO: log error
		// TODO: errors
		return nil, errors.Wrap(err, "failed to collect rows")

		//return nil, fmt.Errorf("failed to collect rows, %s", err.Error())
	}

	if len(apiIntegrations) == 0 {
		return nil, errors.Wrap(err, "no rows for this integration")
	}

	return apiIntegrationToModelSlice(apiIntegrations), nil
}
