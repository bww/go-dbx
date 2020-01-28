package test

type Config struct {
	Name       string
	Migrations string
}

type Option func(c *Config) (*Config, error)

func WithMigrations(rc string) Option {
	return func(c *Config) (*Config, error) {
		c.Migrations = rc
		return c, nil
	}
}
