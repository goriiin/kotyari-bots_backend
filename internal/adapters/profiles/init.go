package profiles

import profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"

type GrpcValidator struct {
	client profiles.ProfilesServiceClient
}

func NewGrpcValidator(client profiles.ProfilesServiceClient) *GrpcValidator {
	return &GrpcValidator{client: client}
}
