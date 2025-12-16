package bots

import (
	"context"

	"github.com/google/uuid"
	bot_grpc "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
)

type Usecase interface {
	Get(ctx context.Context, id uuid.UUID) (model.Bot, error)
}

type Server struct {
	bot_grpc.UnimplementedBotServiceServer
	usecase Usecase
	log     *logger.Logger
}

func NewServer(usecase Usecase, log *logger.Logger) *Server {
	return &Server{
		usecase: usecase,
		log:     log,
	}
}
