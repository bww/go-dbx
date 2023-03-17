package query

type WriteConfig struct {
	Cascade bool
	Params  Params
}

func NewWriteConfig(opts []WriteOption) WriteConfig {
	return WriteConfig{}.WithOptions(opts)
}

func (c WriteConfig) WithOptions(opts []WriteOption) WriteConfig {
	for _, f := range opts {
		c = f(c)
	}
	return c
}

type WriteOption func(WriteConfig) WriteConfig

func UseWriteConfig(c WriteConfig) WriteOption {
	return func(_ WriteConfig) WriteConfig {
		return c
	}
}

func WithCascade() WriteOption {
	return func(c WriteConfig) WriteConfig {
		c.Cascade = true
		return c
	}
}

func WithWriteParams(p Params) WriteOption {
	return func(c WriteConfig) WriteConfig {
		if c.Params == nil {
			c.Params = make(Params)
		}
		for k, v := range p {
			c.Params[k] = v
		}
		return c
	}
}
