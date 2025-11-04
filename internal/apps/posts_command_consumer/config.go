package posts_command_consumer

import (
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_consumer_client"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type PostsCommandConsumerConfig struct {
	config.ConfigBase
	GRPCServerCfg posts_consumer_client.PostsConsGRPCClientConfig `mapstructure:"posts_cons_grpc"`
	Database      postgres.Config                                 `mapstructure:"posts_cons_db"`
	KafkaCons     kafka.KafkaConfig                               `mapstructure:"posts_cons_c"`
	KafkaProd     kafka.KafkaConfig                               `mapstructure:"posts_cons_p"`
}
