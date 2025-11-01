package posts_command_producer

func (p *PostsCommandProducerApp) Close() error {
	return p.producer.Close()
}
