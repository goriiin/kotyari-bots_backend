package api_integrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kotyari-bots_backend/internal/model"
)

func (a *APIIntegrationsRepo) GetIntegrations(ctx context.Context, integrationName string) ([]model.APIIntegration, error) {
	// TODO: log

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
			return nil, fmt.Errorf("no rows for %s", integrationName)
		}
		// TODO: log err
		return nil, fmt.Errorf("unexpected error: %s", err.Error())
	}

	apiIntegrations, err := pgx.CollectRows(rows, pgx.RowToStructByName[APIIntegrationDTO])
	if err != nil {
		// TODO: log error
		// TODO: errors
		return nil, fmt.Errorf("failed to collect rows, %s", err.Error())
	}

	return apiIntegrationToModelSlice(apiIntegrations), nil
}
