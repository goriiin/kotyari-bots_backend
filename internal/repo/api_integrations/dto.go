package api_integrations

import "github.com/kotyari-bots_backend/internal/model"

type APIIntegrationDTO struct {
	Provider string `db:"provider"`
	Url      string `db:"url"`
}

func (a *APIIntegrationDTO) ToModel() model.APIIntegration {
	return model.APIIntegration{
		Provider: a.Provider,
		Url:      a.Url,
	}
}

func apiIntegrationToModelSlice(apiIntegrations []APIIntegrationDTO) []model.APIIntegration {
	apiIntegrationModels := make([]model.APIIntegration, 0, len(apiIntegrations))

	for _, integration := range apiIntegrations {
		apiIntegrationModels = append(apiIntegrationModels, integration.ToModel())
	}

	return apiIntegrationModels
}
