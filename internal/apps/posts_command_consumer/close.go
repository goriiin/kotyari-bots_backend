package posts_command_consumer

func (p *PostsCommandConsumer) Close() error {
	return p.consumer.Close()
}
