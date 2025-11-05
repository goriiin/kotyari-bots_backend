package posts_command_consumer

import (
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_consumer_client"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type PostsCommandConsumerConfig struct {
	config.ConfigBase
	GRPCServerCfg posts_consumer_client.PostsConsGRPCClientConfig `mapstructure:"posts_consumer_grpc"`
	Database      postgres.Config                                 `mapstructure:"posts_database"`
	KafkaCons     kafka.KafkaConfig                               `mapstructure:"posts_consumer_request"`
	KafkaProd     kafka.KafkaConfig                               `mapstructure:"posts_consumer_reply"`
}
