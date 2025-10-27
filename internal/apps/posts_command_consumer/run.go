package posts_command_consumer

func (p *PostsCommandConsumer) Run() error {
	return p.consumerRunner.HandleCommands()
}
