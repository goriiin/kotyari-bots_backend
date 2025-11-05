package posts_command_producer

import (
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/posts_producer_client"
	"github.com/goriiin/kotyari-bots_backend/internal/kafka"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
)

type configAPI struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type PostsCommandProducerConfig struct {
	config.ConfigBase
	API           configAPI                                       `mapstructure:"posts_producer_api"`
	GRPCServerCfg posts_producer_client.PostsProdGRPCClientConfig `mapstructure:"posts_producer_grpc"`
	KafkaProd     kafka.KafkaConfig                               `mapstructure:"posts_producer_request"`
	KafkaCons     kafka.KafkaConfig                               `mapstructure:"posts_producer_reply"`
}
