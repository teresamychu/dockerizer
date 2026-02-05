package commands

type GenerateCommand struct{}

func (c *GenerateCommand) Help() string {
	return ""
}

func (c *GenerateCommand) Synopsis() string {
	return "Generate Dockerfile and docker-compose.yml"
}

func (c *GenerateCommand) Run(args []string) int {
	return 0
}
