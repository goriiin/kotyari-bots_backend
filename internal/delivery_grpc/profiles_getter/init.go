package profiles_getter

import profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"

type ProfileGateway struct {
	client profiles.ProfilesServiceClient
}

func NewProfileGateway(client profiles.ProfilesServiceClient) *ProfileGateway {
	return &ProfileGateway{client: client}
}
