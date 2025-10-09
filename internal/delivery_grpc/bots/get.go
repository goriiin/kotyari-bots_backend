package bots

import (
	"context"

	"github.com/google/uuid"
	botgrpc "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	"github.com/goriiin/kotyari-bots_backend/pkg/constants"
	"github.com/goriiin/kotyari-bots_backend/pkg/ierrors"
)

func (s *Server) GetBot(ctx context.Context, req *botgrpc.GetBotRequest) (*botgrpc.Bot, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, ierrors.DomainToGRPCError(constants.ErrInvalid)
	}

	botModel, err := s.usecase.Get(ctx, id)
	if err != nil {
		return nil, ierrors.DomainToGRPCError(err)
	}

	return &botgrpc.Bot{
		Id:        botModel.ID.String(),
		BotPrompt: botModel.SystemPrompt,
	}, nil
}
