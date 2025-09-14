package profiles

import (
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

func modelToProto(p *model.Profile) *profiles.Profile {
	if p == nil {
		return nil
	}
	return &profiles.Profile{
		Id:     p.ID.String(),
		Name:   p.Name,
		Email:  p.Email,
		Prompt: p.SystemPromt,
	}
}
