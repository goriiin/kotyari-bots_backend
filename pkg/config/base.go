package config

type ConfigBase struct {
	Environment string `mapstructure:"-"`
}

func (c ConfigBase) GetEnvironment() string {
	return c.Environment
}

func (c ConfigBase) IsProduction() bool {
	return c.Environment == envProduction
}

func (c ConfigBase) IsLocal() bool {
	return c.Environment == envLocal
}

type Validator interface {
	Validate() error
}
